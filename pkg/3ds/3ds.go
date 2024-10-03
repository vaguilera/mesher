package f3ds

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/vaguilera/mesher/pkg/vec3"
	"io"
	"os"
)

const (
	MAIN_CHUNK       = 0x4D4D
	EDIT3DS_CHUNK    = 0x3D3D
	EDIT_MATERIAL    = 0xAFFF
	MAT_NAME01       = 0xA000
	KEYF3DS_CHUNK    = 0xB000
	EDITOBJECT_CHUNK = 0x4000
	OBJTRIMESH_CHUNK = 0x4100
	TRIVERTEXL_CHUNK = 0x4110
	TRILOCAL_CHUNK   = 0x4160
	TRIFACEL_CHUNK   = 0x4120
	TRI_MAPPINGCOORS = 0x4140

	COLOR24           = 0x0011
	MATERIAL_AMBIENT  = 0xa010
	MATERIAL_DIFFUSE  = 0xa020
	MATERIAL_SPECULAR = 0xa030
	MATERIAL_TEXMAP   = 0xa200
	MATERIAL_MAPNAME  = 0xa300
)

type F3DS struct {
	Meshes    []Mesh
	Materials []Material
}

type Mesh struct {
	Name        string
	VertexList  []vec3.Vec3
	FacesList   []Face
	CoordsList  []Coord
	NormalsList []vec3.Vec3
}

type Material struct {
	Name     string
	File     string
	Ambient  Color24
	Diffuse  Color24
	Specular Color24
}

type Chunk struct {
	ID   uint16
	Len  uint32
	Data []byte
}

type Face struct {
	V1, V2, V3, Info uint16
}

type Coord struct {
	U, V float32
}

type Color24 struct {
	R, G, B byte
}

func readChunk(file io.Reader) (*Chunk, error) {
	var id uint16
	var length uint32

	err := binary.Read(file, binary.LittleEndian, &id)
	if err != nil {
		return nil, err
	}

	err = binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length-6)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return &Chunk{ID: id, Len: length, Data: data}, nil
}

func find0Byte(data []byte) int {
	idx := 0
	var v byte
	for idx, v = range data {
		if v == 0 {
			break
		}
	}
	return idx
}

func processMaterialTexMap(data []byte) (string, error) {
	r := bytes.NewReader(data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return "", errors.New("no material filename found")
			}
			return "", fmt.Errorf("error reading file: %s", err)
		}
		switch chunk.ID {
		case MATERIAL_MAPNAME:
			idx := find0Byte(chunk.Data)
			return string(chunk.Data[:idx]), nil
		default:
			// fmt.Printf("Ignored chuck: %x\n", chunk.ID)
		}
	}
}

func readColor24(data []byte) (Color24, error) {
	r := bytes.NewReader(data)
	chunk, _ := readChunk(r)
	if chunk.ID != COLOR24 {
		return Color24{}, fmt.Errorf("expecting Color24 chunk. Found: %d", chunk.ID)
	}
	r = bytes.NewReader(chunk.Data)
	var color Color24
	if err := binary.Read(r, binary.LittleEndian, &color); err != nil {
		return Color24{}, fmt.Errorf("error reading Color24: %w", err)
	}
	return color, nil
}

func (f3ds *F3DS) processMaterial(data []byte) (Material, error) {
	r := bytes.NewReader(data)
	var material Material
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return material, nil
			}
			return material, fmt.Errorf("error reading file: %s", err)
		}
		switch chunk.ID {
		case MAT_NAME01:
			idx := find0Byte(chunk.Data)
			material.Name = string(chunk.Data[:idx])
		case MATERIAL_AMBIENT:
			if color, err := readColor24(chunk.Data); err != nil {
				return material, err
			} else {
				material.Ambient = color
			}
		case MATERIAL_DIFFUSE:
			if color, err := readColor24(chunk.Data); err != nil {
				return material, err
			} else {
				material.Diffuse = color
			}
		case MATERIAL_SPECULAR:
			if color, err := readColor24(chunk.Data); err != nil {
				return material, err
			} else {
				material.Specular = color
			}
		case MATERIAL_TEXMAP:
			if name, err := processMaterialTexMap(chunk.Data); err != nil {
				return material, err
			} else {
				material.File = name
			}
		default:
			// fmt.Printf("ffffIgnored chuck: %x - %s\n", chunk.ID, chunk.Data)
		}
	}

}

func (f3ds *F3DS) processCoordinates(data []byte) ([]Coord, error) {
	r := bytes.NewReader(data)
	var numVertex uint16

	if err := binary.Read(r, binary.LittleEndian, &numVertex); err != nil {
		return nil, err
	}

	var coordList []Coord
	for i := uint16(0); i < numVertex; i++ {
		var coord Coord
		if err := binary.Read(r, binary.LittleEndian, &coord); err != nil {
			return nil, fmt.Errorf("error reading texture coordinates: %w", err)
		}
		coordList = append(coordList, coord)
	}
	return coordList, nil
}

func (f3ds *F3DS) processFaces(data []byte) ([]Face, error) {
	r := bytes.NewReader(data)
	var numFaces uint16

	if err := binary.Read(r, binary.LittleEndian, &numFaces); err != nil {
		return nil, err
	}

	var faceList []Face
	for i := uint16(0); i < numFaces; i++ {
		var face Face
		if err := binary.Read(r, binary.LittleEndian, &face); err != nil {
			return nil, fmt.Errorf("error reading faces: %w", err)
		}
		faceList = append(faceList, face)
	}
	return faceList, nil
}

func (f3ds *F3DS) processTriangleList(data []byte) ([]vec3.Vec3, error) {
	r := bytes.NewReader(data)
	var numVert uint16
	if err := binary.Read(r, binary.LittleEndian, &numVert); err != nil {
		return nil, err
	}

	var vertexList []vec3.Vec3
	for i := uint16(0); i < numVert; i++ {
		var vert vec3.Vec3
		if err := binary.Read(r, binary.LittleEndian, &vert); err != nil {
			return nil, fmt.Errorf("error reading vertices: %w", err)
		}
		vertexList = append(vertexList, vert)
	}
	return vertexList, nil
}

func (f3ds *F3DS) processObjTrimesh(data []byte, mesh *Mesh) error {
	r := bytes.NewReader(data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading file: %s", err)
		}

		switch chunk.ID {
		case TRIVERTEXL_CHUNK:
			if mesh.VertexList, err = f3ds.processTriangleList(chunk.Data); err != nil {
				return err
			}
		case TRIFACEL_CHUNK:
			if mesh.FacesList, err = f3ds.processFaces(chunk.Data); err != nil {
				return err
			}
		case TRI_MAPPINGCOORS:
			if mesh.CoordsList, err = f3ds.processCoordinates(chunk.Data); err != nil {
				return err
			}
		default:
			fmt.Printf("Ignored chuck: %x\n", chunk.ID)
		}
	}
}

func (f3ds *F3DS) processEditObject(data []byte) error {
	idx := find0Byte(data)
	name := string(data[:idx])
	mesh := Mesh{Name: name}
	r := bytes.NewReader(data[idx+1:])
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading file: %s", err)
		}
		if chunk.ID == OBJTRIMESH_CHUNK {
			if err := f3ds.processObjTrimesh(chunk.Data, &mesh); err != nil {
				return err
			}
			f3ds.Meshes = append(f3ds.Meshes, mesh)
		}
	}
}

func (f3ds *F3DS) processEditorChunk(data []byte) error {
	r := bytes.NewReader(data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading file: %s", err)
		}
		switch chunk.ID {
		case EDITOBJECT_CHUNK:
			if err = f3ds.processEditObject(chunk.Data); err != nil {
				return err
			}
		case EDIT_MATERIAL:
			if material, err := f3ds.processMaterial(chunk.Data); err != nil {
				return err
			} else {
				f3ds.Materials = append(f3ds.Materials, material)
			}
		default:
			fmt.Printf("Ignored chuck: %x\n", chunk.ID)
		}
	}
}

func (f3ds *F3DS) processFile(f *os.File) error {
	chunk, err := readChunk(f)
	if err != nil {
		return fmt.Errorf("error reading chunk: %s", err)
	}
	if chunk.ID != MAIN_CHUNK {
		return fmt.Errorf("invalid file. Expected 0x4D4D. Got: %x", chunk.ID)
	}
	r := bytes.NewReader(chunk.Data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return errors.New("cant find edit3ds chunk")
			}
			return fmt.Errorf("error reading file: %s", err)
		}
		if chunk.ID == EDIT3DS_CHUNK {
			if err := f3ds.processEditorChunk(chunk.Data); err != nil {
				return err
			}
			break
		}
	}

	for i := range f3ds.Meshes {
		f3ds.GenerateNormals(i)
	}
	return nil
}

func (f3ds *F3DS) GenerateNormals(idx int) {
	vertices := f3ds.Meshes[idx].VertexList
	faces := f3ds.Meshes[idx].FacesList
	vertexNormals := make([]vec3.Vec3, len(vertices))
	faceNormals := make([]vec3.Vec3, len(faces))

	// Calculate face normals
	for i, face := range faces {
		v1 := vec3.Sub(vertices[face.V2], vertices[face.V1])
		v2 := vec3.Sub(vertices[face.V3], vertices[face.V1])
		faceNormals[i] = vec3.CrossProduct(v1, v2)
	}

	// Accumulate normals for each vertex
	for i, face := range faces {
		vertexNormals[face.V1] = vec3.Add(vertexNormals[face.V1], faceNormals[i])
		vertexNormals[face.V2] = vec3.Add(vertexNormals[face.V2], faceNormals[i])
		vertexNormals[face.V3] = vec3.Add(vertexNormals[face.V3], faceNormals[i])
	}

	// Normalize the vertex normals
	for i := range vertexNormals {
		vertexNormals[i].Normalize()
	}

	f3ds.Meshes[idx].NormalsList = vertexNormals
}

func (f3ds *F3DS) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer f.Close()
	return f3ds.processFile(f)
}
