package gltf2

import (
	"github.com/pkg/errors"
)

type MaterialNormalTextureInfo struct {
	Index      *Texture
	TexCoord   IndexTexCoord
	Scale      float32
	Extensions *Extensions
	Extras     *Extras
}

func (s *MaterialNormalTextureInfo) GetExtension() *Extensions {
	return s.Extensions
}

func (s *MaterialNormalTextureInfo) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecMaterialNormalTextureInfo struct {
	Index      *SpecGLTFID     `json:"index"`    // required, minimum(0)
	TexCoord   *IndexTexCoord  `json:"texCoord"` // default(0), minimum(0)
	Scale      *float32        `json:"scale"`    // default(1.0), minimum(0)
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecMaterialNormalTextureInfo) SpecExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecMaterialNormalTextureInfo) Scheme() string {
	return SCHEME_MATERIAL_NORMAL_TEXTUREINFO
}
func (s *SpecMaterialNormalTextureInfo) Syntax(strictness Strictness, root Specifier, parent Specifier) error {

	switch strictness {
	case LEVEL3:
		if s.Scale != nil && *s.Scale < 0.0 {
			return errors.Errorf("MaterialNormalTextureInfo.Scale min(0.0)")
		}
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Index == nil {
			return errors.Errorf("MaterialNormalTextureInfo.Index required")
		}
	}
	return nil
}
func (s *SpecMaterialNormalTextureInfo) To(ctx *parserContext) interface{} {
	res := new(MaterialNormalTextureInfo)
	if s.TexCoord == nil {
		res.TexCoord = TexCoord0
	} else {
		res.TexCoord = *s.TexCoord
	}
	if s.Scale == nil {
		res.Scale = 1.0
	} else {
		res.Scale = *s.Scale
	}
	res.Extras = s.Extras
	return res
}
func (s *SpecMaterialNormalTextureInfo) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.Index, len(Root.Textures)) {
		return errors.Errorf("MaterialNormalTextureInfo.Index linking fail")
	}
	dst.(*MaterialNormalTextureInfo).Index = Root.Textures[*s.Index]
	return nil
}
