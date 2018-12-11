package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

type MaterialPBRMetallicRoughness struct {
	BaseColorFactor          mgl32.Vec4   `json:"baseColorFactor"` // default([1.0, 1.0, 1.0, 1.0]), fixedItem(4), validate(f32Color4)
	MetallicFactor           float32      `json:"metallicFactor"`  // default(1.0), range(0.0, 1.0)
	RoughnessFactor          float32      `json:"roughnessFactor"` // default(1.0), range(0.0, 1.0)
	BaseColorTexture         *TextureInfo `json:"baseColorTexture"`
	MetallicRoughnessTexture *TextureInfo `json:"metallicRoughnessTexture"`
	Extensions               *Extensions  `json:"extensions,omitempty"`
	Extras                   *Extras      `json:"extras,omitempty"`
}

type SpecMaterialPBRMetallicRoughness struct {
	BaseColorFactor          *mgl32.Vec4      `json:"baseColorFactor"` // default([1.0, 1.0, 1.0, 1.0]), fixedItem(4), eachItemRange(0.0, 0.1)
	MetallicFactor           *float32         `json:"metallicFactor"`  // default(1.0), range(0.0, 1.0)
	RoughnessFactor          *float32         `json:"roughnessFactor"` // default(1.0), range(0.0, 1.0)
	BaseColorTexture         *SpecTextureInfo `json:"baseColorTexture"`
	MetallicRoughnessTexture *SpecTextureInfo `json:"metallicRoughnessTexture"`
	Extensions               *Extensions      `json:"extensions,omitempty"`
	Extras                   *Extras          `json:"extras,omitempty"`
}

func (s *SpecMaterialPBRMetallicRoughness) GetChild(i int) ToGLTF {
	switch i {
	case 0:
		return s.BaseColorTexture
	case 1:
		return s.MetallicRoughnessTexture
	}
	return nil
}
func (s *SpecMaterialPBRMetallicRoughness) SetChild(i int, dst, object interface{}) {
	switch i {
	case 0:
		dst.(*MaterialPBRMetallicRoughness).BaseColorTexture = object.(*TextureInfo)
	case 1:
		dst.(*MaterialPBRMetallicRoughness).MetallicRoughnessTexture = object.(*TextureInfo)
	}
}
func (s *SpecMaterialPBRMetallicRoughness) LenChild() int {
	return 2
}

func (s *SpecMaterialPBRMetallicRoughness) Scheme() string {
	return SCHEME_MATERIAL_PBR_METALLIC_ROUGHNESS
}
func (s *SpecMaterialPBRMetallicRoughness) Syntax(strictness Strictness, root interface{}) error {

	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		if s.BaseColorFactor != nil && isValidF32Color4(*s.BaseColorFactor) {
			return errors.WithMessage(ErrorGLTFSpec, "MaterialPBRMetallicRoughness.BaseColorFactor validate(f32Color4)")
		}
		if s.MetallicFactor != nil && (*s.MetallicFactor < 0.0 || *s.MetallicFactor > 1.0) {
			return errors.WithMessage(ErrorGLTFSpec, "MaterialPBRMetallicRoughness.MetallicFactor range(0.0, 1.0)")
		}
		if s.RoughnessFactor != nil && (*s.RoughnessFactor < 0.0 || *s.RoughnessFactor > 1.0) {
			return errors.WithMessage(ErrorGLTFSpec, "MaterialPBRMetallicRoughness.RoughnessFactor range(0.0, 1.0)")
		}
		fallthrough
	case LEVEL1:
	}
	return nil
}
func (s *SpecMaterialPBRMetallicRoughness) To(ctx *parserContext) interface{} {
	res := new(MaterialPBRMetallicRoughness)
	if s.BaseColorFactor == nil {
		res.BaseColorFactor = mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
	} else {
		res.BaseColorFactor = *s.BaseColorFactor
	}
	if s.MetallicFactor == nil {
		res.MetallicFactor = 1.0
	} else {
		res.MetallicFactor = *s.MetallicFactor
	}
	if s.RoughnessFactor == nil {
		res.RoughnessFactor = 1.0
	} else {
		res.RoughnessFactor = *s.RoughnessFactor
	}
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
