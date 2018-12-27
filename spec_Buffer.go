package gltf2

import (
	"github.com/iamGreedy/essence/req"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"path/filepath"
)

// https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#buffers-and-buffer-views
// Implementation Note : Limit of ByteLength
type Buffer struct {
	// nullable
	cache []byte
	unavailableURI bool
	//
	URI *URI
	// If it was nil, it automatically assume size
	ByteLength *int
	Name       string
	Extensions *Extensions
	Extras     *Extras
	// None spec
	UserData interface{}
}

func (s *Buffer) GetExtension() *Extensions {
	return s.Extensions
}

func (s *Buffer) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

func (s *Buffer) Load(storeCache bool) (bts []byte, err error) {
	if s.IsCached() {
		return s.Cache(), nil
	}
	if s.unavailableURI{
		return nil, errors.New("Unavailable URI flag setup")
	}
	// setup 'bts'
	if s.URI == nil {
		// check zero filled buffer
		if s.ByteLength == nil {
			return nil, errors.New("URI nil")
		}
		bts = make([]byte, *s.ByteLength)
	} else {
		rdc, err := req.Standard.Request(s.URI.Data())
		if err != nil {
			return nil, err
		}
		defer rdc.Close()
		bts, err = ioutil.ReadAll(rdc)
		if err != nil {
			return nil, err
		}
	}
	// length limit
	if s.ByteLength != nil {
		bts = bts[:*s.ByteLength]
	}
	// cache
	if storeCache {
		// setup cache
		s.cache = bts
	}
	return bts, nil
}
func (s *Buffer) Modify() ([]byte, error) {
	bts, err := s.Load(true)
	if err != nil {
		return nil, err
	}
	s.unavailableURI = true
	return bts, nil
}
func (s *Buffer) Cache() []byte {
	return s.cache
}
func (s *Buffer) ThrowCache() {
	s.cache = nil
}
func (s *Buffer) IsCached() bool {
	return s.cache != nil
}

type SpecBuffer struct {
	URI        *URI            `json:"URI, omitempty"`
	ByteLength *int            `json:"ByteLength"` // required, min(1)
	Name       *string         `json:"name,omitempty"`
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecBuffer) SpecExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecBuffer) Scheme() string {
	return SCHEME_BUFFER
}
func (s *SpecBuffer) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		if *s.ByteLength < 1 {
			return errors.New("Buffer.ByteLength min(1)")
		}
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.ByteLength == nil {
			return errors.New("Buffer.ByteLength required")
		}
	}
	return nil
}
func (s *SpecBuffer) To(ctx *parserContext) interface{} {
	res := new(Buffer)
	res.URI = s.URI
	if res.URI != nil {
		switch res.URI.Scheme {
		case "":
			fallthrough
		case "file":
			dir := ctx.Directory()
			if dir == "" {
				dir = "."
			}
			res.URI.Path = filepath.Join(dir, filepath.FromSlash(path.Clean("/"+res.URI.Path)))
		}
	}
	res.ByteLength = s.ByteLength
	if s.Name != nil {
		res.Name = *s.Name
	}
	res.Extras = s.Extras
	return res
}
