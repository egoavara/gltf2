package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"image"
)

type OrthographicCamera struct {
	Xmag       float32
	Ymag       float32
	Znear      float32
	Zfar       float32
	Extensions *Extensions
	Extras     *Extras
}

func (s *OrthographicCamera) View(monitorSize image.Point) mgl32.Mat4 {
	return mgl32.Ortho(-s.Xmag, s.Xmag, -s.Ymag, s.Ymag, s.Znear, s.Zfar)
}
func (s *OrthographicCamera) CameraType() CameraType {
	return Orthographic
}

type SpecCameraOrthographic struct {
	Xmag       *float32    `json:"xmag"`  // required, not(0.0)
	Ymag       *float32    `json:"ymag"`  // required, not(0.0)
	Znear      *float32    `json:"znear"` // required, minimum(0.0)
	Zfar       *float32    `json:"zfar"`  // required, larger(0.0), larger(znear)
	Extensions *Extensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`
}

func (s *SpecCameraOrthographic) Scheme() string {
	return SCHEME_CAMERA_ORTHOGRAPHIC
}
func (s *SpecCameraOrthographic) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		if s.Xmag != nil && *s.Xmag == 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Xmag not(0.0)")
		}
		if s.Ymag != nil && *s.Ymag == 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Ymag not(0.0)")
		}

		if s.Znear != nil && *s.Znear < 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Znear minimum(0.0)")
		}
		if s.Zfar != nil && *s.Zfar <= 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Zfar larger(0.0)")
		}
		if s.Znear != nil && s.Zfar != nil && *s.Znear > *s.Zfar {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Zfar larger(znear)")
		}
		fallthrough
	case LEVEL1:
		if s.Xmag == nil {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Xmag required")
		}
		if s.Ymag == nil {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Ymag required")
		}
		if s.Znear == nil {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Znear required")
		}
		if s.Zfar == nil {
			return errors.WithMessage(ErrorGLTFSpec, "OrthographicCamera.Zfar required")
		}
	}
	return nil
}
func (s *SpecCameraOrthographic) To(ctx *parserContext) interface{} {
	res := new(OrthographicCamera)
	res.Xmag = *s.Xmag
	res.Ymag = *s.Ymag
	res.Znear = *s.Znear
	res.Zfar = *s.Zfar
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
