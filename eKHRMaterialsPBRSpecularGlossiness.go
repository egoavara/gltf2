package gltf2

import "github.com/go-gl/mathgl/mgl32"

type KHRMaterialsPBRSpecularGlossiness struct {
	DiffuseFactor             mgl32.Vec4   `json:"diffuseFactor"` // default:[1.0,1.0,1.0,1.0]
	DiffuseTexture            *TextureInfo `json:"diffuseTexture"`
	SpecularFactor            mgl32.Vec3   `json:"specularFactor"`   // default:[1.0,1.0,1.0]
	GlossinessFactor          float32      `json:"glossinessFactor"` // default:1.0
	SpecularGlossinessTexture *TextureInfo `json:"specularGlossinessTexture"`
}

type SpecKHRMaterialsPBRSpecularGlossiness struct {
	DiffuseFactor             *mgl32.Vec4      `json:"diffuseFactor"` // default:[1.0,1.0,1.0,1.0]
	DiffuseTexture            *SpecTextureInfo `json:"diffuseTexture"`
	SpecularFactor            *mgl32.Vec3      `json:"specularFactor"`   // default:[1.0,1.0,1.0]
	GlossinessFactor          *float32         `json:"glossinessFactor"` // default:1.0
	SpecularGlossinessTexture *SpecTextureInfo `json:"specularGlossinessTexture"`
}

func (s *SpecKHRMaterialsPBRSpecularGlossiness) Scheme() string {
	return SCHEME_EXTENSION
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) Syntax(strictness Strictness, root interface{}) error {
	return nil
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) To(ctx *parserContext) interface{} {
	res := new(KHRMaterialsPBRSpecularGlossiness)
	if s.DiffuseFactor == nil {
		res.DiffuseFactor = mgl32.Vec4{1, 1, 1, 1}
	} else {
		res.DiffuseFactor = *s.DiffuseFactor
	}
	if s.SpecularFactor == nil {
		res.SpecularFactor = mgl32.Vec3{1, 1, 1}
	} else {
		res.SpecularFactor = *s.SpecularFactor
	}
	if s.GlossinessFactor == nil {
		res.GlossinessFactor = 1
	} else {
		res.GlossinessFactor = *s.GlossinessFactor
	}
	return res
}

func (s *SpecKHRMaterialsPBRSpecularGlossiness) GetChild(i int) Specifier {
	return s.Children()[i]
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) SetChild(i int, dst, object interface{}) {
	switch s.LenChild() {
	case 2:
		switch i {
		case 0:
			dst.(*KHRMaterialsPBRSpecularGlossiness).DiffuseTexture = object.(*TextureInfo)
		case 1:
			dst.(*KHRMaterialsPBRSpecularGlossiness).SpecularGlossinessTexture = object.(*TextureInfo)
		}
	case 1:
		if s.DiffuseTexture != nil {
			dst.(*KHRMaterialsPBRSpecularGlossiness).DiffuseTexture = object.(*TextureInfo)
		} else {
			dst.(*KHRMaterialsPBRSpecularGlossiness).SpecularGlossinessTexture = object.(*TextureInfo)
		}
	}
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) LenChild() int {
	return len(s.Children())
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) Children() (res []Specifier) {
	if s.DiffuseTexture != nil {
		res = append(res, s.DiffuseTexture)
	}
	if s.SpecularGlossinessTexture != nil {
		res = append(res, s.SpecularGlossinessTexture)
	}
	return res
}
func (s *SpecKHRMaterialsPBRSpecularGlossiness) ImpleGetChild(i int, dst interface{}) interface{} {
	switch s.LenChild() {
	case 2:
		switch i {
		case 0:
			return dst.(*KHRMaterialsPBRSpecularGlossiness).DiffuseTexture
		case 1:
			return dst.(*KHRMaterialsPBRSpecularGlossiness).SpecularGlossinessTexture
		}
	case 1:
		if s.DiffuseTexture != nil {
			return dst.(*KHRMaterialsPBRSpecularGlossiness).DiffuseTexture
		} else {
			return dst.(*KHRMaterialsPBRSpecularGlossiness).SpecularGlossinessTexture
		}
	}
	return nil
}
