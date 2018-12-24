package gltf2

import (
	"github.com/pkg/errors"
)

type Mesh struct {
	Primitives []*MeshPrimitive
	Weights    []float32
	Name       string
	Extensions *Extensions
	Extras     *Extras

	// None spec
	UserData interface{}
}

func (s *Mesh) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecMesh struct {
	Primitives []SpecMeshPrimitive `json:"primitives"` // required, minItem(1)
	Weights    []float32           `json:"weights"`    // minItem(1)
	Name       *string             `json:"name,omitempty"`
	Extensions *SpecExtensions     `json:"extensions,omitempty"`
	Extras     *Extras             `json:"extras,omitempty"`
}

func (s *SpecMesh) GetExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecMesh) Scheme() string {
	return SCHEME_MESH
}
func (s *SpecMesh) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:

		fallthrough
	case LEVEL1:

		if len(s.Primitives) < 0 {
			return errors.Errorf("SpecMeshPrimitive.Primitives required")
		}
	}
	return nil
}
func (s *SpecMesh) To(ctx *parserContext) interface{} {
	res := new(Mesh)
	res.Extras = s.Extras
	res.Primitives = make([]*MeshPrimitive, len(s.Primitives))
	res.Weights = s.Weights
	if s.Name != nil {
		res.Name = *s.Name
	}
	return res
}
func (s *SpecMesh) GetChild(i int) Specifier {
	return &s.Primitives[i]
}
func (s *SpecMesh) SetChild(i int, dst, object interface{}) {
	dst.(*Mesh).Primitives[i] = object.(*MeshPrimitive)
}
func (s *SpecMesh) LenChild() int {
	return len(s.Primitives)
}
func (s *SpecMesh) ImpleGetChild(i int, dst interface{}) interface{} {
	return dst.(*Mesh).Primitives[i]
}
