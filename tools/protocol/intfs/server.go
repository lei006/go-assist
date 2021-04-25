package intfs

type Server interface {
	Start() (bool, error)
	Stop()
}

type ServerManage interface {
	Add(Server)
	Remove(Server)
}
