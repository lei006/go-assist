package daemontool

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lei006/go-assist/daemontool/zcli"
)

type DaemonTool struct {
}

var DefDaemonTool *DaemonTool

func init() {
	DefDaemonTool = &DaemonTool{}
}

func (this *DaemonTool) GetWordPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func (this *DaemonTool) Run(name string, desc string, fn func()) {
	zcli.LaunchServiceRun(name, desc, fn)
}
