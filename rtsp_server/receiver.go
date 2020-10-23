package rtsp_server

import (
	"fmt"
	

	"github.com/deepch/vdk/format/rtsp"
)

type Receiver struct {
	core.PublishContext

	mid    string
	url    string
	client *rtsp.Client
}

func MakeNewReceiver(mid string, url string) *Receiver {

	fmt.Println("name = 创建接收者" + mid)

	recevier := &Receiver{}
	recevier.mid = mid
	recevier.url = url
	//recevier.url = "rtsp://127.0.0.1/test.sdp";
	//recevier.pubsuber = core.NewPublisher("aaaaaaaa",100*time.Millisecond, 10);
	//recevier.pubsuber = core.CreateNewPublisher("", "rtsp", 100, 10)

	go recevier.Run()

	return recevier
}

func (this *Receiver) Run() (err error) {

	return nil
}
