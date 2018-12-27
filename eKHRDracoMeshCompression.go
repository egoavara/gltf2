package gltf2

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// if there is no None compressed mesh, extension Required have KHR_draco_mesh_compression
type KHRDracoMeshCompression struct {
	BufferView *BufferView
	Attributes map[AttributeKey]SpecGLTFID
}
func (s *KHRDracoMeshCompression) ExtensionName() string {
	return "KHR_draco_mesh_compression"
}
func (s *KHRDracoMeshCompression) Constructor(src []byte) (Specifier, error) {
	res := new(SpecKHRDracoMeshCompression)
	if err := json.Unmarshal(src, res); err != nil{
		return nil, err
	}
	return res, nil
}

type SpecKHRDracoMeshCompression struct {
	BufferView *SpecGLTFID                 `json:"bufferView"` // required
	Attributes *map[AttributeKey]SpecGLTFID `json:"attributes"` // required
}

func (s *SpecKHRDracoMeshCompression) Scheme() string {
	return SCHEME_EXTENSION
}
func (s *SpecKHRDracoMeshCompression) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		if !(parent.(*SpecMeshPrimitive).Mode != nil&& (*parent.(*SpecMeshPrimitive).Mode == TRIANGLES || *parent.(*SpecMeshPrimitive).Mode == TRIANGLE_STRIP)){
			return errors.Errorf("KHRDracoMeshCompression >> MeshPrimitive.Mode is TRIANGLES or TRIANGLE_STRIP")
		}
		fallthrough
	case LEVEL1:
		if s.BufferView == nil {
			return errors.Errorf("KHRDracoMeshCompression.BufferView required")
		}
		if s.Attributes == nil {
			return errors.Errorf("KHRDracoMeshCompression.Attributes required")
		}
	}
	return nil
}
func (s *SpecKHRDracoMeshCompression) To(ctx *parserContext) interface{} {
	res := new(KHRDracoMeshCompression)
	res.Attributes = *s.Attributes
	return res
}
func (s *SpecKHRDracoMeshCompression) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.BufferView, len(Root.BufferViews)) {
		return errors.Errorf("KHRDracoMeshCompression.BufferView linking fail, %s", *s.BufferView)
	}
	dst.(*KHRDracoMeshCompression).BufferView = Root.BufferViews[*s.BufferView]
	return nil
}
