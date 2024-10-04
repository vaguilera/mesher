package w3d

import (
	"encoding/base64"
	"encoding/binary"
	f3ds "github.com/vaguilera/mesher/pkg/3ds"
	"math"
	"slices"
)

type W3D struct {
	Meshes []Mesh `json:"meshes"`
}

type Mesh struct {
	Name    string `json:"name"`
	Normals bool   `json:"normals"`
	Coords  bool   `json:"coords"`
	Data    string `json:"data"`
}

func getMeshFrom3DS(mesh f3ds.Mesh) Mesh {
	findexes := make([]byte, len(mesh.FacesList)*6)

	for i, f := range mesh.FacesList {
		idx := i * 6
		binary.LittleEndian.PutUint16(findexes[idx:], f.V1)
		binary.LittleEndian.PutUint16(findexes[idx+2:], f.V2)
		binary.LittleEndian.PutUint16(findexes[idx+4:], f.V3)
	}

	size := 12 // 3 floats
	normals := mesh.NormalsList
	coords := mesh.CoordsList
	isCoords, isNorms := false, false

	if len(normals) > 0 {
		isNorms = true
		size += 12
	}
	if len(coords) > 0 {
		isCoords = true
		size += 8
	}

	vertexData := make([]byte, len(mesh.VertexList)*size)

	for i, v := range mesh.VertexList {
		idx := i * size
		u := math.Float32bits(v.X)
		binary.LittleEndian.PutUint32(vertexData[idx:], u)
		u = math.Float32bits(v.Y)
		binary.LittleEndian.PutUint32(vertexData[idx+4:], u)
		u = math.Float32bits(v.Z)
		binary.LittleEndian.PutUint32(vertexData[idx+8:], u)

		if isNorms {
			idx += 12
			u = math.Float32bits(normals[i].X)
			binary.LittleEndian.PutUint32(vertexData[idx:], u)
			u = math.Float32bits(normals[i].Y)
			binary.LittleEndian.PutUint32(vertexData[idx+4:], u)
			u = math.Float32bits(normals[i].Z)
			binary.LittleEndian.PutUint32(vertexData[idx+8:], u)
		}

		if isCoords {
			idx += 12
			u = math.Float32bits(coords[i].U)
			binary.LittleEndian.PutUint32(vertexData[idx:], u)
			u = math.Float32bits(coords[i].V)
			binary.LittleEndian.PutUint32(vertexData[idx+4:], u)
		}
	}

	data := slices.Concat(findexes, vertexData)
	encodedData := base64.StdEncoding.EncodeToString([]byte(data))

	return Mesh{
		Name:    mesh.Name,
		Normals: isNorms,
		Coords:  isCoords,
		Data:    encodedData,
	}
}

func New3WDFrom3DS(f3ds *f3ds.F3DS) *W3D {
	w3d := W3D{Meshes: []Mesh{}}

	for _, m := range f3ds.Meshes {
		cMesh := getMeshFrom3DS(m)
		w3d.Meshes = append(w3d.Meshes, cMesh)
	}
	return &w3d
}
