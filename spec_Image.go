package gltf2

import (
	"github.com/iamGreedy/essence/req"
	"github.com/pkg/errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"path"
	"path/filepath"
)

type Image interface {
	Name() string
	Extensions() *Extensions
	Extras() *Extras

	Load(useCache bool) (img *image.RGBA, err error)
	Cache() *image.RGBA
	ThrowCache()
	IsCached() bool

	UserData() interface{}
	SetUserData(data interface{})
}

type URIImage struct {
	// nullable
	cache *image.RGBA
	//
	URI *URI
	//
	name       string
	extensions *Extensions
	extras     *Extras
	// None spec
	userData interface{}
}

func (s *URIImage) SetExtension(extensions *Extensions) {
	s.extensions = extensions
}

func (s *URIImage) UserData() interface{} {
	return s.userData
}

func (s *URIImage) SetUserData(data interface{}) {
	s.userData = data
}
func (s *URIImage) Name() string {
	return s.name
}
func (s *URIImage) Extensions() *Extensions {
	return s.extensions
}
func (s *URIImage) Extras() *Extras {
	return s.extras
}
func (s *URIImage) Load(useCache bool) (img *image.RGBA, err error) {
	if s.IsCached() {
		return s.Cache(), nil
	}
	// setup 'img'
	rdc, err := req.Standard.Request(s.URI.Data())
	if err != nil {
		return nil, err
	}
	defer rdc.Close()
	// image decode
	temp, _, err := image.Decode(rdc)
	if err != nil {
		return nil, err
	}
	// image move
	img = image.NewRGBA(temp.Bounds())
	draw.Draw(img, img.Rect, temp, temp.Bounds().Min, draw.Src)
	// cache
	if useCache {
		// setup cache
		s.cache = img
	}
	return img, nil
}
func (s *URIImage) Cache() *image.RGBA {
	return s.cache
}
func (s *URIImage) ThrowCache() {
	s.cache = nil
}
func (s *URIImage) IsCached() bool {
	return s.cache != nil
}

type BufferImage struct {
	// nullable
	cache *image.RGBA
	//
	Mime       MimeType
	BufferView *BufferView
	//
	name       string
	extensions *Extensions
	extras     *Extras
	// None spec
	userData interface{}
}

func (s *BufferImage) SetExtension(extensions *Extensions) {
	s.extensions = extensions
}

func (s *BufferImage) UserData() interface{} {
	return s.userData
}

func (s *BufferImage) SetUserData(data interface{}) {
	s.userData = data
}
func (s *BufferImage) Name() string {
	return s.name
}
func (s *BufferImage) Extensions() *Extensions {
	return s.extensions
}
func (s *BufferImage) Extras() *Extras {
	return s.extras
}
func (s *BufferImage) Load(useCache bool) (img *image.RGBA, err error) {
	rd, err := s.BufferView.LoadReader()
	if err != nil {
		return nil, err
	}
	var temp image.Image
	switch s.Mime {
	case ImagePNG:
		temp, err = png.Decode(rd)
		if err != nil {
			return nil, err
		}
	case ImageJPEG:
		temp, err = jpeg.Decode(rd)
		if err != nil {
			return nil, err
		}
	}
	//
	img = image.NewRGBA(temp.Bounds())
	draw.Draw(img, img.Rect, temp, temp.Bounds().Min, draw.Src)
	//
	return img, nil

}
func (s *BufferImage) Cache() *image.RGBA {
	return s.cache
}
func (s *BufferImage) ThrowCache() {
	s.cache = nil
}
func (s *BufferImage) IsCached() bool {
	return s.cache != nil
}

type SpecImage struct {
	URI        *URI            `json:"URI"`        // exclusive_require(URI, bufferView)
	MimeType   *MimeType       `json:"mimeType"`   //
	BufferView *SpecGLTFID     `json:"bufferView"` // exclusive_require(URI, bufferView), dependency(MimeType)
	Name       *string         `json:"name,omitempty"`
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecImage) GetExtension() *SpecExtensions {
	return s.Extensions
}

func (s *SpecImage) Scheme() string {
	return SCHEME_IMAGE
}
func (s *SpecImage) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if (s.URI != nil && s.BufferView != nil) || (s.URI == nil && s.BufferView == nil) {
			return errors.Errorf("Image must have one of 'Image.URI' or 'Image.bufferView'")
		}
		if s.BufferView != nil && s.MimeType == nil {
			return errors.Errorf("Image.bufferView dependency(MimeType)")
		}
	}
	return nil
}
func (s *SpecImage) To(ctx *parserContext) interface{} {
	if s.URI != nil {
		res := new(URIImage)
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
		if s.Name != nil {
			res.name = *s.Name
		}
		res.extras = s.Extras

		return res
	}
	if s.BufferView != nil {
		res := new(BufferImage)
		//res.BufferView = s.BufferView
		res.Mime = *s.MimeType
		if s.Name != nil {
			res.name = *s.Name
		}
		res.extras = s.Extras

		return res
	}
	panic("Unreachable")
}
func (s *SpecImage) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if bi, ok := dst.(BufferImage); ok {
		if !inRange(*s.BufferView, len(Root.BufferViews)) {
			return errors.Errorf("Image.BufferView linking fail")
		}
		bi.BufferView = Root.BufferViews[*s.BufferView]
	}
	return nil
}
