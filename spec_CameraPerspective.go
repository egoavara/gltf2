package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"image"
	"math"
)

type PerspectiveCamera struct {
	AspectRatio *float32 // if nil, follow monitor ratio
	Yfov        float32
	Znear       float32
	Zfar        float32
	Extensions  *Extensions `json:"extensions,omitempty"`
	Extras      *Extras     `json:"extras,omitempty"`

	// None spec
	userData interface{}
}

func (s *PerspectiveCamera) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

func (s *PerspectiveCamera) GetExtension() *Extensions {
	return s.Extensions
}

func (s *PerspectiveCamera) UserData() interface{} {
	return s.userData
}
func (s *PerspectiveCamera) SetUserData(data interface{}) {
	s.userData = data
}
func (s *PerspectiveCamera) View(monitorSize image.Point) mgl32.Mat4 {
	if s.AspectRatio == nil {
		return perspectiveInfinitable(s.Yfov, float32(monitorSize.X)/float32(monitorSize.Y), s.Znear, s.Zfar)
	}
	return perspectiveInfinitable(s.Yfov, *s.AspectRatio, s.Znear, s.Zfar)
}
func (s *PerspectiveCamera) CameraType() CameraType {
	return Perspective
}

type SpecCameraPerspective struct {
	AspectRatio *float32        `json:"aspectRatio"` // larger(0.0)
	Yfov        *float32        `json:"yfov"`        // required
	Znear       *float32        `json:"znear"`       // required, larger(0.0)
	Zfar        *float32        `json:"zfar"`        // larger(0.0) larger(znear) : not spec but need
	Extensions  *SpecExtensions `json:"extensions,omitempty"`
	Extras      *Extras         `json:"extras,omitempty"`
}

func (s *SpecCameraPerspective) SpecExtension() *SpecExtensions {
	return s.Extensions
}

func (s *SpecCameraPerspective) Scheme() string {
	return SCHEME_CAMERA_PERSPECTIVE
}
func (s *SpecCameraPerspective) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		if s.AspectRatio != nil && *s.AspectRatio <= 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.AspectRatio larger(0.0)")
		}
		if s.Zfar != nil && *s.Zfar <= 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.Zfar larger(0.0)")
		}
		if s.Zfar != nil && *s.Znear > *s.Zfar {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.Zfar larger(znear)")
		}
		if s.Znear != nil && *s.Znear <= 0.0 {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.Znear larger(0.0)")
		}
		fallthrough
	case LEVEL1:
		if s.Yfov == nil {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.Yfov required")
		}
		if s.Znear == nil {
			return errors.WithMessage(ErrorGLTFSpec, "CameraPerspective.Znear required")
		}
	}
	return nil
}
func (s *SpecCameraPerspective) To(ctx *parserContext) interface{} {
	res := new(PerspectiveCamera)
	res.AspectRatio = s.AspectRatio
	res.Yfov = *s.Yfov
	res.Znear = *s.Znear
	if s.Zfar == nil {
		res.Zfar = float32(math.Inf(1))
	} else {
		res.Zfar = *s.Zfar
	}
	res.Extras = s.Extras
	return res
}
