package gltf2

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type GLDefine interface {
	GL() int32
}

type ComponentType int32

func (s ComponentType) GL() int32 {
	return int32(s)
}

// gl.h
const (
	BYTE           ComponentType = 5120
	UNSIGNED_BYTE  ComponentType = 5121
	SHORT          ComponentType = 5122
	UNSIGNED_SHORT ComponentType = 5123
	UNSIGNED_INT   ComponentType = 5125
	FLOAT          ComponentType = 5126
)

func (s ComponentType) Size() int {
	switch s {
	case BYTE:
		return 1
	case UNSIGNED_BYTE:
		return 1
	case SHORT:
		return 2
	case UNSIGNED_SHORT:
		return 2
	case UNSIGNED_INT:
		return 4
	case FLOAT:
		return 4
	}
	return 0
}
func (s ComponentType) String() string {
	switch s {
	case BYTE:
		return "BYTE"
	case UNSIGNED_BYTE:
		return "UNSIGNED_BYTE"
	case SHORT:
		return "SHORT"
	case UNSIGNED_SHORT:
		return "UNSIGNED_SHORT"
	case UNSIGNED_INT:
		return "UNSIGNED_INT"
	case FLOAT:
		return "FLOAT"
	}
	return "nil"
}

type AccessorType uint8

const (
	SCALAR AccessorType = iota
	VEC2   AccessorType = iota
	VEC3   AccessorType = iota
	VEC4   AccessorType = iota
	MAT2   AccessorType = iota
	MAT3   AccessorType = iota
	MAT4   AccessorType = iota
)

func (s AccessorType) Count() int {
	switch s {
	case SCALAR:
		return 1
	case VEC2:
		return 2
	case VEC3:
		return 3
	case VEC4:
		return 4
	case MAT2:
		return 4
	case MAT3:
		return 9
	case MAT4:
		return 16
	}
	return 0
}
func (s AccessorType) String() string {
	switch s {
	case SCALAR:
		return "SCALAR"
	case VEC2:
		return "VEC2"
	case VEC3:
		return "VEC3"
	case VEC4:
		return "VEC4"
	case MAT2:
		return "MAT2"
	case MAT3:
		return "MAT3"
	case MAT4:
		return "MAT4"
	}
	return "nil"
}
func (s *AccessorType) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *AccessorType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {

	case "SCALAR":
		*s = SCALAR
	case "VEC2":
		*s = VEC2
	case "VEC3":
		*s = VEC3
	case "VEC4":
		*s = VEC4

	case "MAT2":
		*s = MAT2
	case "MAT3":
		*s = MAT3
	case "MAT4":
		*s = MAT4
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid AccessorType", v))
	}
	return nil

}

type Path uint8

const (
	Translation Path = iota
	Rotation    Path = iota
	Scale       Path = iota
	Weights     Path = iota
)

func (s Path) String() string {
	switch s {
	case Translation:
		return "translation"
	case Rotation:
		return "rotation"
	case Scale:
		return "scale"
	case Weights:
		return "weights"
	}
	return "nil"
}
func (s *Path) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *Path) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {
	case "translation":
		*s = Translation
	case "rotation":
		*s = Rotation
	case "scale":
		*s = Scale
	case "weights":
		*s = Weights
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid Path", v))
	}
	return nil
}

type Interpolation uint8

const (
	LINEAR      Interpolation = iota
	STEP        Interpolation = iota
	CUBICSPLINE Interpolation = iota
)

func (s Interpolation) String() string {
	switch s {
	case LINEAR:
		return "LINEAR"
	case STEP:
		return "STEP"
	case CUBICSPLINE:
		return "CUBICSPLINE"
	}
	return "nil"
}
func (s *Interpolation) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *Interpolation) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {
	case "LINEAR":
		*s = LINEAR
	case "STEP":
		*s = STEP
	case "CUBICSPLINE":
		*s = CUBICSPLINE
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid Interpolation", v))
	}
	return nil
}

type BufferType int32

func (s BufferType) GL() int32 {
	return int32(s)
}

const (
	// gl.h
	ARRAY_BUFFER         BufferType = 34962
	ELEMENT_ARRAY_BUFFER BufferType = 34963
	// None gl.h
	NEED_TO_DEFINE_BUFFER BufferType = 0
)

func (s BufferType) String() string {
	switch s {
	case ARRAY_BUFFER:
		return "ARRAY_BUFFER"

	case ELEMENT_ARRAY_BUFFER:
		return "ELEMENT_ARRAY_BUFFER"
	}
	return "nil"
}

type CameraType uint8

const (
	Orthographic CameraType = iota
	Perspective  CameraType = iota
)

func (s *CameraType) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *CameraType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {
	case "orthographic":
		*s = Orthographic
	case "perspective":
		*s = Perspective
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid CameraType", v))
	}
	return nil
}

func (s CameraType) String() string {
	switch s {
	case Orthographic:
		return "orthographic"
	case Perspective:
		return "perspective"
	}
	return "nil"
}

type MimeType uint8

const (
	ImagePNG  MimeType = iota
	ImageJPEG MimeType = iota
)

func (s *MimeType) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *MimeType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {
	case "image/png":
		*s = ImagePNG
	case "image/jpeg":
		*s = ImageJPEG
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid MimeType", v))
	}
	return nil
}

func (s MimeType) String() string {
	switch s {
	case ImagePNG:
		return "image/png"
	case ImageJPEG:
		return "image/jpeg"
	}
	return "nil"
}

type AlphaMode uint8

const (
	OPAQUE AlphaMode = iota
	MASK   AlphaMode = iota
	BLEND  AlphaMode = iota
)

func (s *AlphaMode) MarshalJSON() ([]byte, error) {
	if s.String() == "nil" {
		return nil, errors.New("AccessorType invalid")
	}
	return json.Marshal(s.String())
}
func (s *AlphaMode) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	switch v {
	case "OPAQUE":
		*s = OPAQUE
	case "MASK":
		*s = MASK
	case "BLEND":
		*s = BLEND
	default:
		return errors.WithMessage(ErrorEnum, fmt.Sprintf("'%s' is invalid AlphaMode", v))
	}
	return nil
}

func (s AlphaMode) String() string {
	switch s {
	case OPAQUE:
		return "OPAQUE"
	case MASK:
		return "MASK"
	case BLEND:
		return "BLEND"
	}
	return "nil"
}

type Mode int32

func (s Mode) GL() int32 {
	switch s {
	case POINTS:
		return 0
	case LINES:
		return 1
	case LINE_LOOP:
		return 2
	case LINE_STRIP:
		return 3
	case TRIANGLES:
		return 4
	case TRIANGLE_STRIP:
		return 5
	case TRIANGLE_FAN:
		return 6
	default:
		return 0
	}
}

const (
	POINTS         Mode = 0
	LINES          Mode = 1
	LINE_LOOP      Mode = 2
	LINE_STRIP     Mode = 3
	TRIANGLES      Mode = 4
	TRIANGLE_STRIP Mode = 5
	TRIANGLE_FAN   Mode = 6
)

func (s Mode) String() string {
	switch s {
	case POINTS:
		return "POINTS"
	case LINES:
		return "LINES"
	case LINE_LOOP:
		return "LINE_LOOP"
	case LINE_STRIP:
		return "LINE_STRIP"
	case TRIANGLES:
		return "TRIANGLES"
	case TRIANGLE_STRIP:
		return "TRIANGLE_STRIP"
	case TRIANGLE_FAN:
		return "TRIANGLE_FAN"
	}
	return "nil"
}

type AttributeKey string

const (
	POSITION   AttributeKey = "POSITION"
	NORMAL     AttributeKey = "NORMAL"
	TANGENT    AttributeKey = "TANGENT"
	TEXCOORD_0 AttributeKey = "TEXCOORD_0"
	TEXCOORD_1 AttributeKey = "TEXCOORD_1"
	COLOR_0    AttributeKey = "COLOR_0"
	JOINTS_0   AttributeKey = "JOINTS_0"
	WEIGHTS_0  AttributeKey = "WEIGHTS_0"
)
func (s AttributeKey) IsCustom() bool {
	return strings.HasPrefix(string(s), "_")
}
//func (s AttributeKey) AssociateType() ([]ComponentType, []AccessorType) {
//
//}

type MagFilter int32

func (s MagFilter) GL() int32 {
	return int32(s)
}

// gl.h
const (
	MAG_NEAREST MagFilter = 9728
	MAG_LINEAR  MagFilter = 9729
)

func (s MagFilter) String() string {
	switch s {
	case MAG_NEAREST:
		return "NEAREST"
	case MAG_LINEAR:
		return "LINEAR"
	}
	return "nil"
}

type MinFilter int32

func (s MinFilter) GL() int32 {
	return int32(s)
}
func (s MinFilter) IsMipmap() bool {
	switch s {
	case MIN_NEAREST_MIPMAP_NEAREST:
		fallthrough
	case MIN_LINEAR_MIPMAP_NEAREST:
		fallthrough
	case MIN_NEAREST_MIPMAP_LINEAR:
		fallthrough
	case MIN_LINEAR_MIPMAP_LINEAR:
		return true
	}
	return false
}

// gl.h
const (
	MIN_NEAREST MinFilter = 9728
	MIN_LINEAR  MinFilter = 9729

	MIN_NEAREST_MIPMAP_NEAREST MinFilter = 9984
	MIN_LINEAR_MIPMAP_NEAREST  MinFilter = 9985
	MIN_NEAREST_MIPMAP_LINEAR  MinFilter = 9986
	MIN_LINEAR_MIPMAP_LINEAR   MinFilter = 9987
)

func (s MinFilter) String() string {
	switch s {
	case MIN_NEAREST:
		return "NEAREST"
	case MIN_LINEAR:
		return "LINEAR"
	case MIN_NEAREST_MIPMAP_NEAREST:
		return "NEAREST_MIPMAP_NEAREST"
	case MIN_LINEAR_MIPMAP_NEAREST:
		return "LINEAR_MIPMAP_NEAREST"
	case MIN_NEAREST_MIPMAP_LINEAR:
		return "NEAREST_MIPMAP_LINEAR"
	case MIN_LINEAR_MIPMAP_LINEAR:
		return "LINEAR_MIPMAP_LINEAR"
	}
	return "nil"
}

type Wrap int32

func (s Wrap) GL() int32 {
	return int32(s)
}

// gl.h
const (
	CLAMP_TO_EDGE   Wrap = 33071
	MIRRORED_REPEAT Wrap = 33648
	REPEAT          Wrap = 10497
)

func (s Wrap) String() string {
	switch s {
	case CLAMP_TO_EDGE:
		return "CLAMP_TO_EDGE"
	case MIRRORED_REPEAT:
		return "MIRRORED_REPEAT"
	case REPEAT:
		return "REPEAT"
	}
	return "nil"
}

type IndexTexCoord int32

const (
	TexCoord0 IndexTexCoord = 0
	TexCoord1 IndexTexCoord = 1
)

func (s IndexTexCoord) String() string {
	switch s {
	case TexCoord0:
		return "TexCoord_0"
	case TexCoord1:
		return "TexCoord_1"
	}
	return "nil"
}
