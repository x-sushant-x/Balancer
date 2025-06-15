package serverPool

import "github.com/x-sushant-x/Balancer/types"

type ServerPool struct {
	Servers []*types.Server
}

func NewServerPool() *ServerPool {
	return &ServerPool{
		Servers: make([]*types.Server, 0),
	}
}

func (svp *ServerPool) AddServer(server *types.Server) error {
	svp.Servers = append(svp.Servers, server)
	return nil
}

func (svp *ServerPool) GetAllServers() []*types.Server {
	return svp.Servers
}
