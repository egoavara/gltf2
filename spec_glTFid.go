package gltf2

import "fmt"

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/glTFid.schema.json
type SpecGLTFID int

// TODO : minimum 0
func (s SpecGLTFID) String() string {
	return fmt.Sprintf("glTFid(%d)", s)
}
