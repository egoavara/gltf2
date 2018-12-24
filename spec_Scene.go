package gltf2

import "github.com/pkg/errors"

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/scene.schema.json
type Scene struct {
	Nodes      []*Node
	Name       string
	Extensions *Extensions
	Extras     *Extras

	// None spec
	UserData interface{}
}

func (s *Scene) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecScene struct {
	Nodes      []SpecGLTFID `json:"nodes"` // unique, minItem(1)
	Name       *string      `json:"name,omitempty"`
	Extensions *SpecExtensions  `json:"extensions,omitempty"`
	Extras     *Extras      `json:"extras,omitempty"`
}

func (s *SpecScene) GetExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecScene) Scheme() string {
	return SCHEME_SCENE
}
func (s *SpecScene) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if is, id := isUniqueGLTFID(s.Nodes...); !is {
			return errors.Errorf("Scene.Nodes unique '%d'", id)
		}
	}
	return nil

}
func (s *SpecScene) To(ctx *parserContext) interface{} {
	res := new(Scene)
	res.Nodes = make([]*Node, len(s.Nodes))
	if s.Name != nil {
		res.Name = *s.Name
	}
	res.Extras = s.Extras
	return res
}
func (s *SpecScene) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	for i, v := range s.Nodes {
		if !inRange(v, len(Root.Nodes)) {
			return errors.Errorf("Scene.Nodes[%d] linking fail", i)
		}
		dst.(*Scene).Nodes[i] = Root.Nodes[v]
	}
	return nil
}
