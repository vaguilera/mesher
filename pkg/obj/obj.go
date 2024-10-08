package obj

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/vaguilera/mesher/pkg/vec3"
	"log"
	"os"
	"strconv"
	"strings"
)

type ObjFile struct {
	Name        string
	VertexList  []vec3.Vec3
	FacesList   []Face
	CoordsList  []Coord
	NormalsList []vec3.Vec3
}

type Face struct {
	V1, V2, V3 uint16
}

type Coord struct {
	U, V float32
}

func (obj *ObjFile) processFace(tokens []string) error {
	if len(tokens) != 3 {
		return errors.New("sorry. Only triangular faces are supported")
	}
	face := Face{}
	for i, vert := range tokens {
		idxs := strings.Split(vert, "/")
		vidx, err := strconv.ParseUint(idxs[0], 10, 16)
		if err != nil {
			return fmt.Errorf("error parsing float: %s", err.Error())
		}
		vidx16 := uint16(vidx)
		switch i {
		case 0:
			face.V1 = vidx16
		case 1:
			face.V2 = vidx16
		case 2:
			face.V3 = vidx16
		}
	}
	obj.FacesList = append(obj.FacesList, face)
	return nil
}

func (obj *ObjFile) processFile(f *os.File) error {
	scanner := bufio.NewScanner(f)
	cline := 1
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if line[0] == '#' {
			continue
		}

		switch parts[0] {
		case "v":
			x, err1 := strconv.ParseFloat(parts[1], 32)
			y, err2 := strconv.ParseFloat(parts[2], 32)
			z, err3 := strconv.ParseFloat(parts[3], 32)
			x32 := float32(x)
			y32 := float32(y)
			z32 := float32(z)

			if err1 != nil || err2 != nil || err3 != nil {
				return fmt.Errorf("v malformed line %d: %s", cline, line)
			}

			obj.VertexList = append(obj.VertexList, vec3.Vec3{X: x32, Y: y32, Z: z32})
		case "vn":
			x, err1 := strconv.ParseFloat(parts[1], 32)
			y, err2 := strconv.ParseFloat(parts[2], 32)
			z, err3 := strconv.ParseFloat(parts[3], 32)
			x32 := float32(x)
			y32 := float32(y)
			z32 := float32(z)

			if err1 != nil || err2 != nil || err3 != nil {
				return fmt.Errorf("vn malformed line %d: %s", cline, line)
			}

			obj.NormalsList = append(obj.VertexList, vec3.Vec3{X: x32, Y: y32, Z: z32})
		case "vt":
			u, err1 := strconv.ParseFloat(parts[1], 32)
			v, err2 := strconv.ParseFloat(parts[2], 32)
			u32 := float32(u)
			v32 := float32(v)

			if err1 != nil || err2 != nil {
				return fmt.Errorf("vt malformed line %d: %s", cline, line)
			}

			obj.CoordsList = append(obj.CoordsList, Coord{U: u32, V: v32})
		case "f":
			if err := obj.processFace(parts[1:]); err != nil {
				return err
			}
		}

		cline++

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (obj *ObjFile) LoadFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer f.Close()
	return obj.processFile(f)
}
