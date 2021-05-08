package impl

import "github.com/laz/dubbo-go/common/constant"

type Serializer interface {
	Marshal(p DubboPackage) ([]byte, error)
	Unmarshal([]byte, *DubboPackage) error
}

func LoadSerializer(p *DubboPackage) error {
	// NOTE: default serialID is S_Hessian
	serialID := p.Header.SerialID
	if serialID == 0 {
		serialID = constant.S_Hessian2
	}
	serializer, err := GetSerializerById(serialID)
	if err != nil {
		panic(err)
	}
	p.SetSerializer(serializer.(Serializer))
	return nil
}
