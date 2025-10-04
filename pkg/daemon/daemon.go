// pkg/daemon.go
package daemon

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/lei006/go-assist/pkg/zlog"
)

// DaemonService 结构体实现了服务接口
type DaemonService struct {
	Task func()       // 要执行的任务函数
	log  *zlog.Logger // 外部传入的日志对象
}

// Start 服务启动逻辑
func (d *DaemonService) Start(s service.Service) error {

	// 启动传入的任务函数
	go d.run() // 启动后台任务
	return nil
}

// 运行任务
func (d *DaemonService) run() {
	// 调用传入的任务函数
	if d.Task != nil {
		d.Task()
	}
}

// Stop 服务停止逻辑
func (d *DaemonService) Stop(s service.Service) error {

	return nil
}

// StartDaemon 用于创建并启动服务，接受外部传入的配置和日志
func StartDaemon(name, displayName, description string, task func(), logger *zlog.Logger) {
	// 创建服务
	daemon := &DaemonService{
		Task: task,   // 将传入的任务函数赋值给守护进程
		log:  logger, // 传入日志对象
	}

	// 配置服务
	DaemonConfig := &service.Config{
		Name:        name,        // 外部传入的服务名称
		DisplayName: displayName, // 外部传入的显示名称
		Description: description, // 外部传入的描述
	}

	// 创建服务对象
	svc, err := service.New(daemon, DaemonConfig)
	if err != nil {
		logger.Error("Error creating service: ", err) // 使用日志记录错误
		log.Fatalf("Error creating service: %v", err)
	}

	// 定义支持的命令集合
	supportedCommands := map[string]bool{
		"install":   true,
		"uninstall": true,
		"start":     true,
		"stop":      true,
		"restart":   true,
	}

	if len(os.Args) > 1 && supportedCommands[os.Args[1]] {
		err = service.Control(svc, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 启动服务
	err = svc.Run()
	if err != nil {
		logger.Error("Error running service: ", err) // 使用日志记录错误
		log.Fatalf("Error running service: %v", err)
	}
}
