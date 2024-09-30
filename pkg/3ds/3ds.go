package f3ds

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	MAIN_CHUNK       = 0x4D4D
	EDIT3DS_CHUNK    = 0x3D3D
	KEYF3DS_CHUNK    = 0xB000
	EDITOBJECT_CHUNK = 0x4000
	OBJTRIMESH_CHUNK = 0x4100
	TRIVERTEXL_CHUNK = 0x4110
	TRILOCAL_CHUNK   = 0x4160
	TRIFACEL_CHUNK   = 0x4120
)

type F3DS struct {
	vertexList []Vertex
	facesList  []Face
}

type Chunk struct {
	ID   uint16
	Len  uint32
	Data []byte
}

type Vertex struct {
	X, Y, Z float32
}

type Face struct {
	V1, V2, V3, Info uint16
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

func (f3ds *F3DS) processFaces(data []byte) error {
	r := bytes.NewReader(data)
	var numFaces uint16

	if err := binary.Read(r, binary.LittleEndian, &numFaces); err != nil {
		return err
	}
	fmt.Printf("Num of Faces: %d\n", numFaces)

	for i := uint16(0); i < numFaces; i++ {
		var face Face
		if err := binary.Read(r, binary.LittleEndian, &face); err != nil {
			return fmt.Errorf("error reading faces: %w", err)
		}
		f3ds.facesList = append(f3ds.facesList, face)
	}
	return nil
}

func (f3ds *F3DS) processTriangleList(data []byte) error {
	r := bytes.NewReader(data)
	var numVert uint16
	if err := binary.Read(r, binary.LittleEndian, &numVert); err != nil {
		return err
	}
	fmt.Printf("Num of vertex: %d\n", numVert)

	for i := uint16(0); i < numVert; i++ {
		var vert Vertex
		if err := binary.Read(r, binary.LittleEndian, &vert); err != nil {
			return fmt.Errorf("error reading vertices: %w", err)
		}
		f3ds.vertexList = append(f3ds.vertexList, vert)
	}
	return nil
}

func (f3ds *F3DS) processObjTrimesh(data []byte) error {
	r := bytes.NewReader(data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				if len(f3ds.vertexList) == 0 {
					return errors.New("error. No triangle data found in mesh")
				} else {
					return nil
				}
			}
			return fmt.Errorf("error reading file: %s", err)
		}

		switch chunk.ID {
		case TRIVERTEXL_CHUNK:
			if err = f3ds.processTriangleList(chunk.Data); err != nil {
				return err
			}
		case TRIFACEL_CHUNK:
			if err = f3ds.processFaces(chunk.Data); err != nil {
				return err
			}
		default:
			fmt.Printf("Ignored chuck: %x\n", chunk.ID)
		}
	}
}

func (f3ds *F3DS) processEditObject(data []byte) error {
	idx := 0
	var v byte
	for idx, v = range data {
		if v == 0 {
			break
		}
	}
	name := string(data[:idx])
	fmt.Printf("object name: %s\n", name)
	r := bytes.NewReader(data[idx+1:])
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return errors.New("error. No triangle data found in mesh")
			}
			return fmt.Errorf("error reading file: %s", err)
		}
		if chunk.ID == OBJTRIMESH_CHUNK {
			return f3ds.processObjTrimesh(chunk.Data)
		}
	}
}

func (f3ds *F3DS) processEditorChunk(data []byte) error {
	r := bytes.NewReader(data)
	for {
		chunk, err := readChunk(r)
		if err != nil {
			if err == io.EOF {
				return errors.New("cannot find any mesh in the file")
			}
			return fmt.Errorf("error reading file: %s", err)
		}
		if chunk.ID == EDITOBJECT_CHUNK {
			return f3ds.processEditObject(chunk.Data)
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
			return f3ds.processEditorChunk(chunk.Data)
		}
	}
}

func (f3ds *F3DS) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer f.Close()
	err = f3ds.processFile(f)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
