package handler

import (
	"encoding/binary"
	"errors"
)

import (
	"github.com/apache/dubbo-getty"
)

type PackageHandler struct{}

func (h *PackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	dataLen := len(data)
	if dataLen < 4 {
		return nil, 0, nil
	}

	start := 0
	pos := start + 4
	pkgLen := int(binary.LittleEndian.Uint32(data[start:pos]))
	if dataLen < pos+pkgLen {
		return nil, pos + pkgLen, nil
	}
	start = pos

	pos = start + pkgLen
	s := string(data[start:pos])

	return s, pos, nil
}

func (h *PackageHandler) Write(ss getty.Session, p interface{}) ([]byte, error) {
	pkg, ok := p.(string)
	if !ok {
		Log.Infof("illegal pkg:%+v", p)
		return nil, errors.New("invalid package")
	}

	pkgLen := int32(len(pkg))
	pkgStreams := make([]byte, 0, 4+len(pkg))

	// pkg len
	start := 0
	pos := start + 4
	binary.LittleEndian.PutUint32(pkgStreams[start:pos], uint32(pkgLen))
	start = pos

	// pkg
	pos = start + int(pkgLen)
	copy(pkgStreams[start:pos], pkg[:])

	return pkgStreams[:pos], nil
}
