package manage_server

import "go-assist/protocol/intfs"

type ManageServer struct {
}

func MakeManageServer() intfs.ServerManage {
	srv := &ManageServer{}

	return srv
}

func (this *ManageServer) Add(server intfs.Server) {

}

func (this *ManageServer) Remove(server intfs.Server) {

}
