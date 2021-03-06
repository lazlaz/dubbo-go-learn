package remoting

import "github.com/laz/dubbo-go/common"

// It is interface of server for network communication.
// If you use getty as network communication, you should define GettyServer that implements this interface.
type Server interface {
	//invoke once for connection
	Start()
	//it is for destroy
	Stop()
}

// This is abstraction level. it is like facade.
type ExchangeServer struct {
	Server Server
	Url    *common.URL
}

// Create ExchangeServer
func NewExchangeServer(url *common.URL, server Server) *ExchangeServer {
	exchangServer := &ExchangeServer{
		Server: server,
		Url:    url,
	}
	return exchangServer
}

// start server
func (server *ExchangeServer) Start() {
	server.Server.Start()
}
