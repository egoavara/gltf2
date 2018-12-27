package gltf2

type Sampler struct {
	MagFilter MagFilter
	MinFilter MinFilter
	WrapS     Wrap
	WrapT     Wrap

	Name       string
	Extensions *Extensions
	Extras     *Extras

	// None spec
	UserData interface{}
}

func (s *Sampler) GetExtension() *Extensions {
	return s.Extensions
}

func (s *Sampler) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

func DefaultSampler() *Sampler {
	return &Sampler{
		MagFilter: MAG_LINEAR,
		MinFilter: MIN_NEAREST_MIPMAP_LINEAR,
		WrapS:     REPEAT,
		WrapT:     REPEAT,
	}
}

type SpecSampler struct {
	MagFilter *MagFilter `json:"magFilter"` // notspec default(MAG_LINEAR) : [https://www.khronos.org/registry/OpenGL-Refpages/gl4/html/glTexParameter.xhtm]
	MinFilter *MinFilter `json:"minFilter"` // notspec default(MIN_NEAREST_MIPMAP_LINEAR) : [https://www.khronos.org/registry/OpenGL-Refpages/gl4/html/glTexParameter.xhtm]
	WrapS     *Wrap      `json:"wrapS"`     // notspec default(REPEAT) : [https://www.khronos.org/registry/OpenGL-Refpages/gl4/html/glTexParameter.xhtm]
	WrapT     *Wrap      `json:"wrapT"`     // notspec default(REPEAT) : [https://www.khronos.org/registry/OpenGL-Refpages/gl4/html/glTexParameter.xhtm]

	Name       *string         `json:"name,omitempty"`
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
}

func (s *SpecSampler) SpecExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecSampler) Scheme() string {
	return SCHEME_SAMPLER
}
func (s *SpecSampler) Syntax(strictness Strictness, root Specifier, parent Specifier) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
	}
	return nil
}
func (s *SpecSampler) To(ctx *parserContext) interface{} {
	res := new(Sampler)
	if s.MagFilter == nil {
		res.MagFilter = MAG_LINEAR
	} else {
		res.MagFilter = *s.MagFilter
	}
	if s.MinFilter == nil {
		res.MinFilter = MIN_NEAREST_MIPMAP_LINEAR
	} else {
		res.MinFilter = *s.MinFilter
	}
	if s.WrapS == nil {
		res.WrapS = REPEAT
	} else {
		res.WrapS = *s.WrapS
	}
	if s.WrapT == nil {
		res.WrapT = REPEAT
	} else {
		res.WrapT = *s.WrapT
	}

	return res
}
