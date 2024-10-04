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
    Indexes  string `json:"indexes"`
    Vertices string `json:"vertices"`
}
```

W3D is just a JSON file with an object containing an array of Meshes.

Every mesh contains a mesh name, 2 boolean that indicate if the mesh contains normals and texture coordinates and 
the indexes and vertices data encoded in base64. 
As the W3D file has been designed to work with WebGPU, indexes data is encoded as a series of 16 bit words,
and vertices are a series of structs with the following format?

```go
vertices: []float32 {
	X,Y,Z float32 // vertex position 
	NX,NY,NZ float32 // vertex normal
	U,V float32 // texture coordinates
}
```

If the mesh doesn't have normals or texture coordinates, they are simply not included.
Every value is little-endian encoded.

# ToDo
- Add command flag to select meshes to export (currently evey mesh in the 3DS file is included)
- Add support for OBJ files
- Add command flag to skip normals generation