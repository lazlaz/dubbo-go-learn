package getty

import (
	"errors"
	getty "github.com/apache/dubbo-getty"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/laz/dubbo-go/common/logger"
	"github.com/laz/dubbo-go/remoting"
	perrors "github.com/pkg/errors"
	"reflect"
)

// RpcServerPackageHandler Read data from client and Write data to client
type RpcServerPackageHandler struct {
	server *Server
}

func NewRpcServerPackageHandler(server *Server) *RpcServerPackageHandler {
	return &RpcServerPackageHandler{server: server}
}

// Read data from client. if the package size from client is larger than 4096 byte, client will read 4096 byte
// and send to client each time. the Read can assemble it.
func (p *RpcServerPackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	req, length, err := (p.server.codec).Decode(data)
	//resp,len, err := (*p.).DecodeResponse(buf)
	if err != nil {
		if errors.Is(err, hessian.ErrHeaderNotEnough) || errors.Is(err, hessian.ErrBodyNotEnough) {
			return nil, 0, nil
		}

		logger.Errorf("pkg.Unmarshal(ss:%+v, len(@data):%d) = error:%+v", ss, len(data), err)

		return nil, 0, err
	}

	return req, length, err
}

// Write send the data to client
func (p *RpcServerPackageHandler) Write(ss getty.Session, pkg interface{}) ([]byte, error) {
	res, ok := pkg.(*remoting.Response)
	if ok {
		buf, err := (p.server.codec).EncodeResponse(res)
		if err != nil {
			logger.Warnf("binary.Write(res{%#v}) = err{%#v}", res, perrors.WithStack(err))
			return nil, perrors.WithStack(err)
		}
		return buf.Bytes(), nil
	}

	req, ok := pkg.(*remoting.Request)
	if ok {
		buf, err := (p.server.codec).EncodeRequest(req)
		if err != nil {
			logger.Warnf("binary.Write(req{%#v}) = err{%#v}", res, perrors.WithStack(err))
			return nil, perrors.WithStack(err)
		}
		return buf.Bytes(), nil
	}

	logger.Errorf("illegal pkg:%+v\n, it is %+v", pkg, reflect.TypeOf(pkg))
	return nil, perrors.New("invalid rpc response")

}
