package gltf2

import (
	"github.com/pkg/errors"
)

type MeshPrimitive struct {
	Attributes map[AttributeKey]*Accessor
	Indices    *Accessor
	Material   *Material
	Mode       Mode
	Targets    []map[AttributeKey]*Accessor
	Extensions *Extensions
	Extras     *Extras

	// None spec
	UserData interface{}
}

type SpecMeshPrimitive struct {
	Attributes map[AttributeKey]SpecGLTFID   `json:"attributes"` // required, minItem(1)
	Indices    *SpecGLTFID                   `json:"indices"`    //
	Material   *SpecGLTFID                   `json:"material"`   //
	Mode       *Mode                         `json:"mode"`       // default(TRIANGLES)
	Targets    []map[AttributeKey]SpecGLTFID `json:"targets"`    // [*]allow(POSITION, NORMAL, TANGENT)
	Extensions *Extensions                   `json:"extensions,omitempty"`
	Extras     *Extras                       `json:"extras,omitempty"`
}

func (s *SpecMeshPrimitive) Scheme() string {
	return SCHEME_MESH_PRIMITIVE
}
func (s *SpecMeshPrimitive) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if len(s.Attributes) < 0 {
			return errors.Errorf("MeshPrimitive.Attributes required")
		}
		for _, target := range s.Targets {
			count := 0
			if _, ok := target[POSITION]; ok {
				count++
			}
			if _, ok := target[NORMAL]; ok {
				count++
			}
			if _, ok := target[TANGENT]; ok {
				count++
			}
			if len(target) != count {
				return errors.Errorf("MeshPrimitive.Targets [*]allow(POSITION, NORMAL, TANGENT)")
			}
		}
	}
	return nil
}
func (s *SpecMeshPrimitive) To(ctx *parserContext) interface{} {
	res := new(MeshPrimitive)
	if s.Attributes != nil {
		res.Attributes = make(map[AttributeKey]*Accessor)
	}
	if s.Targets != nil {
		res.Targets = make([]map[AttributeKey]*Accessor, len(s.Targets))
		for i := range res.Targets {
			res.Targets[i] = make(map[AttributeKey]*Accessor)
		}
	}
	if s.Mode == nil {
		res.Mode = TRIANGLES
	} else {
		res.Mode = *s.Mode
	}
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
func (s *SpecMeshPrimitive) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	for k, v := range s.Attributes {
		if !inRange(v, len(Root.Accessors)) {
			return errors.Errorf("MeshPrimitive.Attributes[%s] linking fail", k)
		}
		dst.(*MeshPrimitive).Attributes[k] = Root.Accessors[v]
	}
	if s.Indices != nil {
		if !inRange(*s.Indices, len(Root.Accessors)) {
			return errors.Errorf("MeshPrimitive.Indices linking fail")
		}
		dst.(*MeshPrimitive).Indices = Root.Accessors[*s.Indices]
	}
	if s.Material != nil {
		if !inRange(*s.Material, len(Root.Materials)) {
			return errors.Errorf("MeshPrimitive.Material linking fail")
		}
		dst.(*MeshPrimitive).Material = Root.Materials[*s.Material]
	}

	for i, target := range s.Targets {
		for k, v := range target {
			if !inRange(v, len(Root.Accessors)) {
				return errors.Errorf("MeshPrimitive.Target[%d][%s] linking fail", i, k)
			}
			dst.(*MeshPrimitive).Targets[i][k] = Root.Accessors[v]
		}
	}
	return nil
}
