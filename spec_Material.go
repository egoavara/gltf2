package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

type Material struct {
	Name                 string                        `json:"name,omitempty"`
	Extensions           *Extensions                   `json:"extensions,omitempty"`
	Extras               *Extras                       `json:"extras,omitempty"`
	PBRMetallicRoughness *MaterialPBRMetallicRoughness `json:"pbrMetallicRoughness"`
	NormalTexture        *MaterialNormalTextureInfo    `json:"normalTexture"`
	OcclusionTexture     *MaterialOcclusionTextureInfo `json:"occlusionTexture"`
	EmissiveFactor       mgl32.Vec3                    `json:"emissiveFactor"`  // default([0.0, 0.0, 0.0]
	EmissiveTexture      *TextureInfo                  `json:"emissiveTexture"` //
	AlphaMode            AlphaMode                     `json:"alphaMode"`       // default(OPAQUE)
	AlphaCutoff          float32                       `json:"alphaCutoff"`     // default(0.5), minimum(0.0), dependency(AlphaMode) ! ignore(dependency) : cause default
	DoubleSided          bool                          `json:"doubleSided"`     // default(false)
}

type SpecMaterial struct {
	Name                 *string                           `json:"name,omitempty"`
	Extensions           *Extensions                       `json:"extensions,omitempty"`
	Extras               *Extras                           `json:"extras,omitempty"`
	PBRMetallicRoughness *SpecMaterialPBRMetallicRoughness `json:"pbrMetallicRoughness"`
	NormalTexture        *SpecMaterialNormalTextureInfo    `json:"normalTexture"`
	OcclusionTexture     *SpecMaterialOcclusionTextureInfo `json:"occlusionTexture"`
	EmissiveFactor       *mgl32.Vec3                       `json:"emissiveFactor"`  // default([0.0, 0.0, 0.0], validate(f32Color)
	EmissiveTexture      *SpecTextureInfo                  `json:"emissiveTexture"` //
	AlphaMode            *AlphaMode                        `json:"alphaMode"`       // default(OPAQUE)
	AlphaCutoff          *float32                          `json:"alphaCutoff"`     // default(0.5), minimum(0.0), dependency(AlphaMode) ! ignore(dependency) : cause default
	DoubleSided          *bool                             `json:"doubleSided"`     // default(false)
}

func (s *SpecMaterial) Scheme() string {
	return SCHEME_MATERIAL
}
func (s *SpecMaterial) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		// it was normally LEVEL1 strictness, but AlphaMode has default so it was LEVEL3
		if s.AlphaCutoff != nil && s.AlphaMode == nil {
			return errors.WithMessage(ErrorGLTFSpec, "Material.AlphaCutoff dependency(AlphaMode)")
		}
		fallthrough
	case LEVEL2:
		if s.EmissiveFactor != nil && isValidF32Color3(*s.EmissiveFactor) {
			return errors.WithMessage(ErrorGLTFSpec, "Material.EmissiveFactor validate(f32Color)")
		}
		if s.AlphaCutoff != nil && *s.AlphaCutoff < 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "Material.AlphaCutoff minimum(0.0)")
		}
		fallthrough
	case LEVEL1:

	}
	return nil
}
func (s *SpecMaterial) To(ctx *parserContext) interface{} {
	res := new(Material)
	if s.Name != nil {
		res.Name = *s.Name
	}
	res.Extensions = s.Extensions
	res.Extras = s.Extras
	if s.EmissiveFactor == nil {
		res.EmissiveFactor = mgl32.Vec3{0, 0, 0}
	} else {
		res.EmissiveFactor = *s.EmissiveFactor
	}
	if s.AlphaMode == nil {
		res.AlphaMode = OPAQUE
	} else {
		res.AlphaMode = *s.AlphaMode
	}
	if s.AlphaCutoff == nil {
		res.AlphaCutoff = .5
	} else {
		res.AlphaCutoff = *s.AlphaCutoff
	}
	if s.DoubleSided == nil {
		res.DoubleSided = false
	} else {
		res.DoubleSided = *s.DoubleSided
	}
	return res
}

func (s *SpecMaterial) GetChild(i int) ToGLTF {
	return s.Childrun()[i]
}

func (s *SpecMaterial) SetChild(i int, dst, object interface{}) {
	if dsto, ok := dst.(*Material); ok {
		switch s.Childrun()[i].(type) {
		case *SpecMaterialPBRMetallicRoughness:
			dsto.PBRMetallicRoughness = object.(*MaterialPBRMetallicRoughness)
		case *SpecMaterialNormalTextureInfo:
			dsto.NormalTexture = object.(*MaterialNormalTextureInfo)
		case *SpecMaterialOcclusionTextureInfo:
			dsto.OcclusionTexture = object.(*MaterialOcclusionTextureInfo)
		case *SpecTextureInfo:
			dsto.EmissiveTexture = object.(*TextureInfo)
		}
	}
}

func (s *SpecMaterial) LenChild() int {
	return len(s.Childrun())
}
func (s *SpecMaterial) Childrun() (res []ToGLTF) {
	if s.PBRMetallicRoughness != nil {
		res = append(res, s.PBRMetallicRoughness)
	}
	if s.NormalTexture != nil {
		res = append(res, s.NormalTexture)
	}
	if s.OcclusionTexture != nil {
		res = append(res, s.OcclusionTexture)
	}
	if s.EmissiveTexture != nil {
		res = append(res, s.OcclusionTexture)
	}
	return res
}
func (s *SpecMaterial) ImpleGetChild(i int, dst interface{}) interface{} {
	switch s.Childrun()[i].(type) {
	case *SpecMaterialPBRMetallicRoughness:
		return dst.(*Material).PBRMetallicRoughness
	case *SpecMaterialNormalTextureInfo:
		return dst.(*Material).NormalTexture
	case *SpecMaterialOcclusionTextureInfo:
		return dst.(*Material).OcclusionTexture
	case *SpecTextureInfo:
		return dst.(*Material).EmissiveTexture
	}
	return nil
}
