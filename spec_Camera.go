package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"image"
)

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/camera.schema.json
// It makes
type Camera interface {
	CameraType() CameraType
	View(monitorSize image.Point) mgl32.Mat4
	UserData() interface{}
	SetUserData(data interface{})
}

type SpecCamera struct {
	Type         *CameraType             `json:"type"`         // required
	Orthographic *SpecCameraOrthographic `json:"orthographic"` // link(type)
	Perspective  *SpecCameraPerspective  `json:"perspective"`  // link(type)
}

func (s *SpecCamera) Scheme() string {
	return SCHEME_CAMERA
}
func (s *SpecCamera) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Type == nil {
			return errors.Errorf("Camera.AccessorType required")
		}
		switch *s.Type {
		case Orthographic:
			if s.Orthographic == nil {
				return errors.Errorf("Camera.Orthographic must need when Camera.AccessorType == Orthographic")
			}
			if err := s.Orthographic.Syntax(strictness, nil); err != nil {
				return err
			}
		case Perspective:
			if s.Perspective == nil {
				return errors.Errorf("Camera.Perspective must need when Camera.AccessorType == Perspective")
			}
			if err := s.Perspective.Syntax(strictness, nil); err != nil {
				return err
			}
		default:
			// TODO : enum constraint
			return errors.Errorf("Camera.AccessorType should be one of [Orthographic, Perspective] ")
		}
	}

	return nil
}
func (s *SpecCamera) To(ctx *parserContext) interface{} {
	switch *s.Type {
	case Orthographic:
		if s.Orthographic == nil {
			return errors.Errorf("Camera.Orthographic must need when Camera.AccessorType == Orthographic")
		}
		return s.Orthographic.To(nil)
	case Perspective:
		if s.Perspective == nil {
			return errors.Errorf("Camera.Perspective must need when Camera.AccessorType == Perspective")
		}
		return s.Perspective.To(nil)
	}
	panic("Unreachable")
}
