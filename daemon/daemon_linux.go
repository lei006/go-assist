package daemon

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
)

type (
	linuxSystemService struct {
		name        string
		detect      func() bool
		interactive func() bool
		new         func(i Ife, c *Config) (ServiceIfe, error)
	}
	systemd struct {
		i Ife
		*Config
	}
)

const (
	optionReloadSignal = "ReloadSignal"
	optionPIDFile      = "PIDFile"
)

var errNoUserServiceSystemd = errors.New("user services are not supported on systemd")

func init() {
	chooseSystem(linuxSystemService{
		name:   "linux-systemd",
		detect: isSystemd,
		interactive: func() bool {
			is, _ := isInteractive()
			return is
		},
		new: newSystemdService,
	})
}

func (sc linuxSystemService) String() string {
	return sc.name
}

func (sc linuxSystemService) Detect() bool {
	return sc.detect()
}

func (sc linuxSystemService) Interactive() bool {
	return sc.interactive()
}

func (sc linuxSystemService) New(i Ife, c *Config) (s ServiceIfe, err error) {
	s, err = sc.new(i, c)
	if err == nil {
		err = isSudo()
	}
	return
}

func isInteractive() (bool, error) {
	return os.Getppid() != 1, nil
}

var tf = map[string]interface{}{
	"cmd": func(s string) string {
		return `"` + strings.Replace(s, `"`, `\"`, -1) + `"`
	},
	"cmdEscape": func(s string) string {
		return strings.Replace(s, " ", `\x20`, -1)
	},
}

func isSystemd() bool {
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return true
	}
	return false
}

func newSystemdService(i Ife, c *Config) (ServiceIfe, error) {
	s := &systemd{
		i:      i,
		Config: c,
	}

	return s, nil
}

func (s *systemd) String() string {
	if len(s.DisplayName) > 0 {
		return s.DisplayName
	}
	return s.Name
}

func (s *systemd) configPath() (cp string, err error) {
	if s.Option.Bool(optionUserService, optionUserServiceDefault) {
		err = errNoUserServiceSystemd
		return
	}
	cp = "/etc/systemd/system/" + s.Config.Name + ".service"
	return
}
func (s *systemd) template() *template.Template {
	return template.Must(template.New("").Funcs(tf).Parse(systemdScript))
}

func (s *systemd) Install() error {
	confPath, err := s.configPath()
	if err != nil {
		return err
	}
	_, err = os.Stat(confPath)
	if err == nil {
		return fmt.Errorf("init already exists: %s", confPath)
	}

	f, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer f.Close()

	path := s.execPath()
	var to = &struct {
		*Config
		Path         string
		ReloadSignal string
		PIDFile      string
	}{
		s.Config,
		path,
		s.Option.String(optionReloadSignal, ""),
		s.Option.String(optionPIDFile, ""),
	}

	err = s.template().Execute(f, to)
	if err != nil {
		return err
	}

	err = run("systemctl", "enable", s.Name+".service")
	if err != nil {
		return err
	}
	return run("systemctl", "daemon-reload")
}

func (s *systemd) Uninstall() error {
	_ = run("systemctl", "stop", s.Name+".service")
	err := run("systemctl", "disable", s.Name+".service")
	if err != nil {
		return err
	}
	cp, err := s.configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(cp); err != nil {
		return err
	}
	return nil
}

func (s *systemd) Run() (err error) {
	err = s.i.Start(s)
	if err != nil {
		return err
	}

	s.Option.FuncSingle(optionRunWait, func() {
		var sigChan = make(chan os.Signal, 3)
		signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigChan
	})()

	return s.i.Stop(s)
}

func (s *systemd) Start() error {
	if os.Getuid() == 0 {
		return run("systemctl", "start", s.Name+".service")
	} else {
		return run("sudo", "-n", "systemctl", "start", s.Name+".service")
	}
}

func (s *systemd) Stop() error {
	if os.Getuid() == 0 {
		return run("systemctl", "stop", s.Name+".service")
	} else {
		return run("sudo", "-n", "systemctl", "stop", s.Name+".service")
	}
}

func (s *systemd) Restart() error {
	if os.Getuid() == 0 {
		return run("systemctl", "restart", s.Name+".service")
	} else {
		return run("sudo", "-n", "systemctl", "restart", s.Name+".service")
	}
}

func (s *systemd) Status() string {
	var res string
	if os.Getuid() == 0 {
		res, _ = runGrep("running", "systemctl", "status", s.Name+".service")
	} else {
		res, _ = runGrep("running", "sudo", "-n", "systemctl", "status", s.Name+".service")
	}
	if res != "" {
		return "Running"
	}
	return "Stop"
}

const systemdScript = `[Unit]
Description={{.Description}}
After=network.target

[Service]
StartLimitInterval=5
StartLimitBurst=10
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .RootDir}}RootDirectory={{.RootDir|cmd}}{{end}}
{{if .WorkingDir}}WorkingDirectory={{.WorkingDir|cmdEscape}}{{end}}
{{if .UserName}}User={{.UserName}}{{end}}
{{if .ReloadSignal}}ExecReload=/bin/kill -{{.ReloadSignal}} "$MAINPID"{{end}}
{{if .PIDFile}}PIDFile={{.PIDFile|cmd}}{{end}}
Restart=always
RestartSec=120ms
EnvironmentFile=-/etc/sysconfig/{{.Name}}
KillMode=process
TimeoutStopSec=1s

[Install]
WantedBy=multi-user.target
`
