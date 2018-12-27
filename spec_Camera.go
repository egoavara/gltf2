package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"image"
)

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/camera.schema.json
type Camera struct {
	Setting    CameraSetting
	Extensions *Extensions
	Extras     *Extras
}

func (s *Camera) GetExtension() *Extensions {
	return s.Extensions
}

func (s *Camera) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type CameraSetting interface {
	CameraType() CameraType
	View(monitorSize image.Point) mgl32.Mat4
	UserData() interface{}
	SetUserData(data interface{})
	GetExtension() *Extensions
}

type SpecCamera struct {
	Type         *CameraType             `json:"type"`         // required
	Orthographic *SpecCameraOrthographic `json:"orthographic"` // link(type)
	Perspective  *SpecCameraPerspective  `json:"perspective"`  // link(type)

	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecCamera) SpecExtension() *SpecExtensions {
	return s.Extensions
}

func (s *SpecCamera) GetChild(i int) Specifier {
	switch *s.Type {
	case Orthographic:
		return s.Orthographic
	case Perspective:
		return s.Perspective
	}
	return nil
}
func (s *SpecCamera) SetChild(i int, dst, object interface{}) {
	switch *s.Type {
	case Orthographic:
		dst.(*Camera).Setting = object.(CameraSetting)
	case Perspective:
		dst.(*Camera).Setting = object.(CameraSetting)
	}
}
func (s *SpecCamera) LenChild() int {
	return 1
}
func (s *SpecCamera) ImpleGetChild(i int, dst interface{}) interface{} {
	return dst.(*Camera).Setting
}

func (s *SpecCamera) Scheme() string {
	return SCHEME_CAMERA
}
func (s *SpecCamera) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Type == nil {
			return errors.Errorf("CameraSetting.AccessorType required")
		}
	}
	return nil
}
func (s *SpecCamera) To(ctx *parserContext) interface{} {
	//switch *s.Type {
	//case Orthographic:
	//	if s.Orthographic == nil {
	//		return errors.Errorf("CameraSetting.Orthographic must need when CameraSetting.AccessorType == Orthographic")
	//	}
	//	return s.Orthographic.To(nil)
	//case Perspective:
	//	if s.Perspective == nil {
	//		return errors.Errorf("CameraSetting.Perspective must need when CameraSetting.AccessorType == Perspective")
	//	}
	//	return s.Perspective.To(nil)
	//}
	//panic("Unreachable")
	res := new(Camera)
	res.Extras = s.Extras
	return res
}
