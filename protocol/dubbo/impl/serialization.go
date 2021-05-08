package impl

import (
	"fmt"
	"github.com/laz/dubbo-go/common/constant"
)

var (
	serializers = make(map[string]interface{})
	nameMaps    = make(map[byte]string)
)

func init() {
	nameMaps = map[byte]string{
		constant.S_Hessian2: constant.HESSIAN2_SERIALIZATION,
		constant.S_Proto:    constant.PROTOBUF_SERIALIZATION,
	}
}

func SetSerializer(name string, serializer interface{}) {
	serializers[name] = serializer
}
func GetSerializerById(id byte) (interface{}, error) {
	name, ok := nameMaps[id]
	if !ok {
		panic(fmt.Sprintf("serialId %d not found", id))
	}
	serializer, ok := serializers[name]
	if !ok {
		panic(fmt.Sprintf("serialization %s not found", name))
	}
	return serializer, nil
}
