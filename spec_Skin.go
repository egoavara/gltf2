package gltf2

import (
	"github.com/pkg/errors"
)

type Skin struct {
	InverseBindMatrices *Accessor
	Skeleton            *Node
	Joints              []*Node
	Name                string
	Extensions          *Extensions
	Extras              *Extras

	// None spec
	UserData interface{}
}

func (s *Skin) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

func (s *Skin) GetExtension() *Extensions {
	return s.Extensions
}

type SpecSkin struct {
	InverseBindMatrices *SpecGLTFID     `json:"inverseBindMatrices"` // When undefined, it is Ident4x4 matrix
	Skeleton            *SpecGLTFID     `json:"skeleton"`            // When undefined, joints transforms resolve to scene root.
	Joints              []SpecGLTFID    `json:"joints"`              // require(min = 1), unique
	Name                *string         `json:"name,omitempty"`
	Extensions          *SpecExtensions `json:"extensions,omitempty"`
	Extras              *Extras         `json:"extras,omitempty"`
}

func (s *SpecSkin) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if s.InverseBindMatrices != nil{
		if !inRange(*s.InverseBindMatrices, len(Root.Accessors)){
			return errors.Errorf("Skin.InverseBindMatrices linking fail")
		}
		dst.(*Skin).InverseBindMatrices = Root.Accessors[*s.InverseBindMatrices]
	}
	if s.Skeleton != nil{
		if !inRange(*s.Skeleton, len(Root.Nodes)){
			return errors.Errorf("Skin.Skeleton linking fail")
		}
		dst.(*Skin).Skeleton = Root.Nodes[*s.Skeleton]
	}
	if len(s.Joints) > 0{
		dst.(*Skin).Joints = make([]*Node, 0, len(s.Joints))
		for i, v := range s.Joints {
			if !inRange(v, len(Root.Nodes)){
				return errors.Errorf("Skin.Joints linking fail %d", i)
			}
			dst.(*Skin).Joints = append(dst.(*Skin).Joints, Root.Nodes[v])
		}
	}
	return nil
}

func (s *SpecSkin) Scheme() string {
	return SCHEME_SKIN
}

func (s *SpecSkin) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if len(s.Joints) == 0 {
			return errors.Errorf("Skin.Joints require(min = 1)")
		}
		if s.Joints != nil {
			if ok, id := isUniqueGLTFID(s.Joints...); !ok {
				return errors.Errorf("Skin.Joints unique, but id %d overlap", id)
			}
		}
	}
	return nil
}

func (s *SpecSkin) To(ctx *parserContext) interface{} {
	res := new(Skin)
	if s.Name != nil{
		res.Name = *s.Name
	}
	res.Extras = s.Extras
	return res
}

func (s *SpecSkin) SpecExtension() *SpecExtensions {
	return s.Extensions
}
