package gltf2

import "github.com/pkg/errors"

type Texture struct {
	Sampler    *Sampler    `json:"sampler"`
	Source     Image       `json:"source"`
	Name       string      `json:"name,omitempty"`
	Extensions *Extensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`

	// None spec
	UserData interface{}
}

func (s *Texture) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecTexture struct {
	Sampler    *SpecGLTFID `json:"sampler"`
	Source     *SpecGLTFID `json:"source"`
	Name       *string     `json:"name,omitempty"`
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`
}

func (s *SpecTexture) GetExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecTexture) Scheme() string {
	return SCHEME_TEXTURE
}
func (s *SpecTexture) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
	}
	return nil
}
func (s *SpecTexture) To(ctx *parserContext) interface{} {
	res := new(Texture)
	if s.Name == nil {
		res.Name = ""
	} else {
		res.Name = *s.Name
	}
	res.Extras = s.Extras
	return res
}
func (s *SpecTexture) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if s.Sampler == nil {
		dst.(*Texture).Sampler = DefaultSampler()
	} else {
		if !inRange(*s.Sampler, len(Root.Samplers)) {
			return errors.Errorf("Texture.Sampler linking fail")
		}
		dst.(*Texture).Sampler = Root.Samplers[*s.Sampler]
	}
	if s.Source == nil {
		// TODO : https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#texture-data
		// 		Unclear, what if Source is undefined?
		//		> maybe default image for Texture source?
		panic("Undefined works")
	} else {
		if !inRange(*s.Source, len(Root.Images)) {
			return errors.Errorf("Texture.Source linking fail")
		}
		dst.(*Texture).Source = Root.Images[*s.Source]
	}
	return nil
}
