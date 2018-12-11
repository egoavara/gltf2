package gltf2

import (
	"github.com/pkg/errors"
)

type MaterialOcclusionTextureInfo struct {
	Index      *Texture
	TexCoord   IndexTexCoord
	Strength   float32
	Extensions *Extensions
	Extras     *Extras
}

type SpecMaterialOcclusionTextureInfo struct {
	Index      *SpecGLTFID    `json:"index"`    // required, minimum(0)
	TexCoord   *IndexTexCoord `json:"texCoord"` // default(0), minimum(0)
	Strength   *float32       `json:"scale"`    // default(1.0), range(0.0, 1.0)
	Extensions *Extensions    `json:"extensions,omitempty"`
	Extras     *Extras        `json:"extras,omitempty"`
}

func (s *SpecMaterialOcclusionTextureInfo) Scheme() string {
	return SCHEME_MATERIAL_OCCLUSION_TEXTUREINFO
}
func (s *SpecMaterialOcclusionTextureInfo) Syntax(strictness Strictness, root interface{}) error {

	switch strictness {
	case LEVEL3:
		if s.Strength != nil && *s.Strength < 0.0 || *s.Strength > 1.0 {
			return errors.WithMessage(ErrorGLTFSpec, "MaterialOcclusionTextureInfo.Strength range(0.0, 1.0)")
		}
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Index == nil {
			return errors.WithMessage(ErrorGLTFSpec, "MaterialOcclusionTextureInfo.Index required")
		}
	}
	return nil
}
func (s *SpecMaterialOcclusionTextureInfo) To(ctx *parserContext) interface{} {
	res := new(MaterialOcclusionTextureInfo)
	if s.TexCoord == nil {
		res.TexCoord = TexCoord0
	} else {
		res.TexCoord = *s.TexCoord
	}
	if s.Strength == nil {
		res.Strength = 0
	} else {
		res.Strength = *s.Strength
	}
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
func (s *SpecMaterialOcclusionTextureInfo) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.Index, len(Root.Textures)) {
		return errors.Errorf("MaterialOcclusionTextureInfo.Index linking fail")
	}
	dst.(*MaterialOcclusionTextureInfo).Index = Root.Textures[*s.Index]
	return nil
}
