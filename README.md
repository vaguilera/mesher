# mesher
Convert 3DS formats into web-usable format

Current version only support Autodesk 3DS files

Usage:
``Mesher <input.3ds> <output file>``

W3D Format:
```go
type W3D struct {
    Meshes []Mesh `json:"meshes"`
}

type Mesh struct {
    Name    string `json:"name"`
    Normals bool   `json:"normals"`
    Coords  bool   `json:"coords"`
    Data    string `json:"data"`
}
```

W3D is just a JSON file with an object containing an array of Meshes.

Every mesh contains a mesh name, 2 boolean that indicate if the mesh contains normals and texture coordinates and the data in base64. 
As the W3D file has been designed to work with WebGPU, data is a bytes buffer with this format:

```go
faceIndexes: []uint16 // 3 floats per triangle
// After the indexes buffer, it comes the vertices info. For each vertex you can expect:
vertices: []float32 {
	X,Y,Z float32 // vertex position 
	NX,NY,NZ float32 // vertex normal
	U,V float32 // texture coordinates
}
// If the mesh doesn't have normals or texture coordinates, they are simply not included
```

# ToDo
- Add command flag to select meshes to export (currently evey mesh in the 3DS file is included)
- Add support for OBJ files
- Add command flag to skip normals generation