package livertc

import (
	cmap "github.com/orcaman/concurrent-map"
)

// 如果要直接修改数据，调用这个...
type EngineState struct {
	pub_cmap cmap.ConcurrentMap //存放定阅者...

}

//发布一个通道
func (state *EngineState) PublishChannel() {
	//http.Serve()

}
