package w3d

import (
	"encoding/binary"
	"fmt"
	f3ds "github.com/vaguilera/mesher/pkg/3ds"
)

type W3D struct {
}

func New3WDFrom3DS(f3ds *f3ds.F3DS) {
	findexes := make([]byte, len(f3ds.Meshes[0].FacesList)*6)

	//buf := [2]byte{}
	for i, f := range f3ds.Meshes[0].FacesList {
		binary.LittleEndian.PutUint16(findexes[i*6:], f.V1)
		binary.LittleEndian.PutUint16(findexes[i*6+2:], f.V2)
		binary.LittleEndian.PutUint16(findexes[i*6+4:], f.V3)
	}

	fmt.Println(findexes)
	fmt.Println(f3ds.Meshes[0].FacesList)
}
