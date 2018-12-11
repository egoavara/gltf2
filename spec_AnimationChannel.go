package gltf2

import (
	"github.com/pkg/errors"
)

type AnimationChannel struct {
	Sampler    *AnimationSampler       `json:"sampler"` // required
	Target     *AnimationChannelTarget `json:"target"`  // required
	Extensions *Extensions             `json:"extensions,omitempty"`
	Extras     *Extras                 `json:"extras,omitempty"`
}

type SpecAnimationChannel struct {
	Sampler    *SpecGLTFID                 `json:"sampler"` // required
	Target     *SpecAnimationChannelTarget `json:"target"`  // required
	Extensions *Extensions                 `json:"extensions,omitempty"`
	Extras     *Extras                     `json:"extras,omitempty"`
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
	if err := s.Target.Syntax(strictness, nil); err != nil {
		return err
	}
	return nil
}
func (s *SpecAnimationChannel) To(ctx *parserContext) interface{} {
	res := new(AnimationChannel)
	res.Target = s.Target.To(nil).(*AnimationChannelTarget)
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
func (s *SpecAnimationChannel) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if !inRange(*s.Sampler, len(parent.(*Animation).Samplers)) {
		return errors.Errorf("AnimationChannel.Sampler linking fail")
	}
	dst.(*AnimationChannel).Sampler = parent.(*Animation).Samplers[*s.Sampler]
	if err := s.Target.Link(Root, dst, dst.(*AnimationChannel).Target); err != nil {
		return err
	}
	return nil
}
