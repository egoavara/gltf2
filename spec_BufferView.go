package gltf2

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
)

// https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#buffers-and-buffer-views
// Implementation Note : BufferView.Target determine BufferType
//
type BufferView struct {
	Buffer     *Buffer // Linking
	ByteOffset int     // default 0
	ByteLength int     // must have
	ByteStride int     // default 0
	Target     BufferType
	Name       string
	Extensions *Extensions
	Extras     *Extras
	// None spec
	UserData interface{}
}

func (s *BufferView) GetExtension() *Extensions {
	return s.Extensions
}

func (s *BufferView) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

func (s *BufferView) Load() ([]byte, error) {

	var bts []byte
	if s.Buffer.IsCached() {
		bts = s.Buffer.Cache()
	} else {
		var err error
		bts, err = s.Buffer.Load(false)
		if err != nil {
			return nil, err
		}
	}
	return bts[s.ByteOffset : s.ByteOffset+s.ByteLength], nil
}
func (s *BufferView) LoadReader() (io.Reader, error) {
	bts, err := s.Load()
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bts), nil
}

type SpecBufferView struct {
	Buffer     *SpecGLTFID     `json:"buffer"`     // required
	ByteOffset *int            `json:"byteOffset"` // default(0), min(0)
	ByteLength *int            `json:"ByteLength"` // required, min(1)
	ByteStride *int            `json:"byteStride"` // range(4, 252, step=4) ! default(0) : not spec, but can be
	Target     *BufferType     `json:"target"`     //
	Name       *string         `json:"name,omitempty"`
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecBufferView) SpecExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecBufferView) Scheme() string {
	return SCHEME_BUFFERVIEW
}
func (s *SpecBufferView) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		if s.ByteStride != nil && *s.ByteStride < 4 && *s.ByteStride > 252 && *s.ByteStride%4 == 0 {
			return errors.Errorf("BufferView.ByteStride range(4, 252, step=4), but got '%d'", *s.ByteStride)
		}
		if s.ByteLength != nil && *s.ByteLength < 1 {
			return errors.Errorf("BufferView.ByteLength min(1), but got '%d'", *s.ByteLength)
		}
		if s.ByteOffset != nil && *s.ByteOffset < 0 {
			return errors.Errorf("BufferView.ByteOffset min(0), but got '%d'", *s.ByteOffset)
		}
		fallthrough
	case LEVEL2:

		fallthrough
	case LEVEL1:
		if s.Buffer == nil {
			return errors.Errorf("BufferView.Buffer required")
		}
		if s.ByteLength == nil {
			return errors.Errorf("BufferView.ByteLength required")
		}
	}
	return nil
}
func (s *SpecBufferView) To(ctx *parserContext) interface{} {
	res := new(BufferView)
	if s.ByteOffset == nil {
		res.ByteOffset = 0
	} else {
		res.ByteOffset = *s.ByteOffset
	}
	if s.ByteLength == nil {
		res.ByteLength = 0
	} else {
		res.ByteLength = *s.ByteLength
	}
	if s.ByteStride == nil {
		res.ByteStride = 0
	} else {
		res.ByteStride = *s.ByteStride
	}
	if s.Target == nil {
		res.Target = NEED_TO_DEFINE_BUFFER
	} else {
		res.Target = *s.Target
	}
	if s.Name != nil {
		res.Name = *s.Name
	}
	res.Extras = s.Extras
	return res
}
func (s *SpecBufferView) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.Buffer, len(Root.Buffers)) {
		return errors.Errorf("BufferView.Buffer linking fail")
	}
	dst.(*BufferView).Buffer = Root.Buffers[*s.Buffer]
	return nil
}
