package common

import (
	"github.com/laz/dubbo-go/common/logger"
	"reflect"
	"strings"
	"sync"
)
import (
	perrors "github.com/pkg/errors"
)

const (
	METHOD_MAPPER = "MethodMapper"
)

// RPCService
// rpc service interface
type RPCService interface {
	// Reference:
	// rpc service id or reference id
	Reference() string
}

// Service is description of service
type Service struct {
	name     string
	rcvr     reflect.Value
	rcvrType reflect.Type
	methods  map[string]*MethodType
}

// MethodType is description of service method.
type MethodType struct {
	method    reflect.Method
	ctxType   reflect.Type   // request context
	argsType  []reflect.Type // args except ctx, include replyType if existing
	replyType reflect.Type   // return value, otherwise it is nil
}
type serviceMap struct {
	mutex        sync.RWMutex                   // protects the serviceMap
	serviceMap   map[string]map[string]*Service // protocol -> service name -> service
	interfaceMap map[string][]*Service          // interface -> service
}

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	// ServiceMap store description of service.
	ServiceMap = &serviceMap{
		serviceMap:   make(map[string]map[string]*Service),
		interfaceMap: make(map[string][]*Service),
	}
)

// Register registers a service by @interfaceName and @protocol
func (sm *serviceMap) Register(interfaceName, protocol, group, version string, rcvr RPCService) (string, error) {
	if sm.serviceMap[protocol] == nil {
		sm.serviceMap[protocol] = make(map[string]*Service)
	}
	if sm.interfaceMap[interfaceName] == nil {
		sm.interfaceMap[interfaceName] = make([]*Service, 0, 16)
	}

	s := new(Service)
	s.rcvrType = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if sname == "" {
		s := "no service name for type " + s.rcvrType.String()
		logger.Errorf(s)
		return "", perrors.New(s)
	}
	/*if !isExported(sname) {
		s := "type " + sname + " is not exported"
		logger.Errorf(s)
		return "", perrors.New(s)
	}*/

	sname = ServiceKey(interfaceName, group, version)
	if server := sm.GetService(protocol, interfaceName, group, version); server != nil {
		return "", perrors.New("service already defined: " + sname)
	}
	s.name = sname
	s.methods = make(map[string]*MethodType)

	// Install the methods
	methods := ""
	methods, s.methods = suitableMethods(s.rcvrType)

	if len(s.methods) == 0 {
		s := "type " + sname + " has no exported methods of suitable type"
		logger.Errorf(s)
		return "", perrors.New(s)
	}
	sm.mutex.Lock()
	sm.serviceMap[protocol][s.name] = s
	sm.interfaceMap[interfaceName] = append(sm.interfaceMap[interfaceName], s)
	sm.mutex.Unlock()

	return strings.TrimSuffix(methods, ","), nil
}

// GetService gets a service definition by protocol and name
func (sm *serviceMap) GetService(protocol, interfaceName, group, version string) *Service {
	serviceKey := ServiceKey(interfaceName, group, version)
	return sm.GetServiceByServiceKey(protocol, serviceKey)
}

// GetService gets a service definition by protocol and service key
func (sm *serviceMap) GetServiceByServiceKey(protocol, serviceKey string) *Service {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	if s, ok := sm.serviceMap[protocol]; ok {
		if srv, ok := s[serviceKey]; ok {
			return srv
		}
		return nil
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ
func suitableMethods(typ reflect.Type) (string, map[string]*MethodType) {
	methods := make(map[string]*MethodType)
	var mts []string
	logger.Debugf("[%s] NumMethod is %d", typ.String(), typ.NumMethod())
	method, ok := typ.MethodByName(METHOD_MAPPER)
	var methodMapper map[string]string
	if ok && method.Type.NumIn() == 1 && method.Type.NumOut() == 1 && method.Type.Out(0).String() == "map[string]string" {
		methodMapper = method.Func.Call([]reflect.Value{reflect.New(typ.Elem())})[0].Interface().(map[string]string)
	}

	for m := 0; m < typ.NumMethod(); m++ {
		method = typ.Method(m)
		if mt := suiteMethod(method); mt != nil {
			methodName, ok := methodMapper[method.Name]
			if !ok {
				methodName = method.Name
			}
			methods[methodName] = mt
			mts = append(mts, methodName)
		}
	}
	return strings.Join(mts, ","), methods
}

// suiteMethod returns a suitable Rpc methodType
func suiteMethod(method reflect.Method) *MethodType {
	mtype := method.Type
	mname := method.Name
	inNum := mtype.NumIn()
	outNum := mtype.NumOut()

	// Method must be exported.
	if method.PkgPath != "" {
		return nil
	}

	var (
		replyType, ctxType reflect.Type
		argsType           []reflect.Type
	)

	// this method is in RPCService
	// we force users must implement RPCService interface in their provider
	// and RPCService has only one method "Reference"
	// In general, this method should not be exported to client
	// so we ignore this method
	// see RPCService
	if mname == "Reference" {
		return nil
	}

	if outNum != 1 && outNum != 2 {
		logger.Warnf("method %s of mtype %v has wrong number of in out parameters %d; needs exactly 1/2",
			mname, mtype.String(), outNum)
		return nil
	}

	// The latest return type of the method must be error.
	if returnType := mtype.Out(outNum - 1); returnType != typeOfError {
		logger.Warnf("the latest return type %s of method %q is not error", returnType, mname)
		return nil
	}

	// replyType
	if outNum == 2 {
		replyType = mtype.Out(0)
		/*if !isExportedOrBuiltinType(replyType) {
			logger.Errorf("reply type of method %s not exported{%v}", mname, replyType)
			return nil
		}*/
	}

	index := 1

	// ctxType
	if inNum > 1 && mtype.In(1).String() == "context.Context" {
		ctxType = mtype.In(1)
		index = 2
	}

	for ; index < inNum; index++ {
		argsType = append(argsType, mtype.In(index))
		/*// need not be a pointer.
		if !isExportedOrBuiltinType(mtype.In(index)) {
			logger.Errorf("argument type of method %q is not exported %v", mname, mtype.In(index))
			return nil
		}*/
	}

	return &MethodType{method: method, argsType: argsType, replyType: replyType, ctxType: ctxType}
}
