package gltf2

type SpecKHRDracoMeshCompression struct {
	BufferView SpecGLTFID `json:"bufferView"`
	Attributes []map[AttributeKey]SpecGLTFID `json:"attributes"`
}