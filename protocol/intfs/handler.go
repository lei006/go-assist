package intfs

import "fmt"

type Info struct {
	Key   string
	URL   string
	UID   string
	Inter bool
}

func (info Info) IsInterval() bool {
	return info.Inter
}

func (info Info) String() string {
	return fmt.Sprintf("<key: %s, URL: %s, UID: %s, Inter: %v>",
		info.Key, info.URL, info.UID, info.Inter)
}

type HandlerPacket interface {
	HandleReader(ReadPacketer)
	HandleWriter(WritePacketer)
}

type ReadPacketer interface {
	Info() Info
	Close(error)
	Alive() bool
	Read(*Packet) error
}

type WritePacketer interface {
	Info() Info
	Close(error)
	Alive() bool
	CalcBaseTimestamp()
	Write(*Packet) error
}
