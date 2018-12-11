package gltf2

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Buffers []Buffer

// https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#buffers-and-buffer-views
// Implementation Note : Limit of ByteLength
type Buffer struct {
	// nullable
	cache []byte
	//
	URI *URI
	// If it was nil, it automatically assume size
	ByteLength *int
	Name       string
	Extensions *Extensions
	Extras     *Extras
}

func (s *Buffer) Load(useCache bool) (bts []byte, err error) {
	if s.IsCached() {
		return s.Cache(), nil
	}
	// setup 'bts'
	if s.URI == nil {
		// check zero filled buffer
		if s.ByteLength == nil {
			return nil, errors.New("URI nil")
		}
		bts = make([]byte, *s.ByteLength)
	} else {
		path := s.URI.Data()
		switch path.Scheme {
		case "http":
			// http server
			fallthrough
		case "https":
			// http TLS server
			var res *http.Response
			res, err = http.Get(path.String())
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			bts, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
		case "":
			fallthrough
		case "file":
			// local file
			var f *os.File
			f, err = os.Open(path.Path)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			bts, err = ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.Errorf("Unsupported scheme '%s'", path.Scheme)
		}
	}
	// length limit
	if s.ByteLength != nil {
		bts = bts[:*s.ByteLength]
	}
	// cache
	if useCache {
		// setup cache
		s.cache = bts
	}
	return bts, nil
}
func (s *Buffer) Cache() []byte {
	return s.cache
}
func (s *Buffer) IsCached() bool {
	return s.cache != nil
}

type SpecBuffer struct {
	URI        *URI        `json:"URI, omitempty"`
	ByteLength *int        `json:"ByteLength"` // required, min(1)
	Name       *string     `json:"name,omitempty"`
	Extensions *Extensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`
}

func (s *SpecBuffer) Scheme() string {
	return SCHEME_BUFFER
}
func (s *SpecBuffer) Syntax(strictness Strictness, root interface{}) error {
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
	res.Extensions = s.Extensions
	res.Extras = s.Extras
	return res
}
