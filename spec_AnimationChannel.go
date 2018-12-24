package gltf2

import (
	"github.com/pkg/errors"
)

type AnimationChannel struct {
	Sampler    *AnimationSampler       `json:"sampler"` // required
	Target     *AnimationChannelTarget `json:"target"`  // required
	Extensions *Extensions             `json:"extensions,omitempty"`
	Extras     *Extras                 `json:"extras,omitempty"`
	// None spec
	UserData interface{}
}

func (s *AnimationChannel) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecAnimationChannel struct {
	Sampler    *SpecGLTFID                 `json:"sampler"` // required
	Target     *SpecAnimationChannelTarget `json:"target"`  // required
	Extensions *SpecExtensions             `json:"extensions,omitempty"`
	Extras     *Extras                     `json:"extras,omitempty"`
}

func (s *SpecAnimationChannel) GetExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecAnimationChannel) Scheme() string {
	return SCHEME_ANIMATION_CHANNEL
}
func (s *SpecAnimationChannel) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Sampler == nil {
			return errors.WithMessage(ErrorGLTFSpec, "AnimationChannel.Sampler required")
		}
		if s.Target == nil {
			return errors.WithMessage(ErrorGLTFSpec, "AnimationChannel.Target required")
		}
	}
	return nil
}
func (s *SpecAnimationChannel) To(ctx *parserContext) interface{} {
	res := new(AnimationChannel)
	res.Extras = s.Extras
	return res
}
func (s *SpecAnimationChannel) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.Sampler, len(parent.(*Animation).Samplers)) {
		return errors.Errorf("AnimationChannel.Sampler linking fail")
	}
	dst.(*AnimationChannel).Sampler = parent.(*Animation).Samplers[*s.Sampler]
	return nil
}

func (s *SpecAnimationChannel) GetChild(i int) Specifier {
	return s.Target
}
func (s *SpecAnimationChannel) SetChild(i int, dst, object interface{}) {
	dst.(*AnimationChannel).Target = object.(*AnimationChannelTarget)
}
func (s *SpecAnimationChannel) LenChild() int {
	if s.Target == nil {
		return 0
	}
	return 1
}
func (s *SpecAnimationChannel) ImpleGetChild(i int, dst interface{}) interface{} {
	return dst.(*AnimationChannel).Target
}
