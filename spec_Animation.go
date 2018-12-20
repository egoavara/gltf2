package gltf2

import (
	"github.com/pkg/errors"
)

type Animation struct {
	Channels   []*AnimationChannel `json:"channels"` // required, minItem(1)
	Samplers   []*AnimationSampler `json:"samplers"` // required, minItem(1)
	Name       string              `json:"name,omitempty"`
	Extensions *Extensions         `json:"extensions,omitempty"`
	Extras     *Extras             `json:"extras,omitempty"`
	// None spec
	UserData interface{}
}

type SpecAnimation struct {
	Channels   []SpecAnimationChannel `json:"channels"` // required, minItem(1)
	Samplers   []SpecAnimationSampler `json:"samplers"` // required, minItem(1)
	Name       *string                `json:"name,omitempty"`
	Extensions *Extensions            `json:"extensions,omitempty"`
	Extras     *Extras                `json:"extras,omitempty"`
}

func (s *SpecAnimation) Scheme() string {
	return SCHEME_ANIMATION
}

func (s *SpecAnimation) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Channels == nil {
			return errors.WithMessage(ErrorGLTFSpec, "NewPlayer.Channels required")
		}
		if s.Samplers == nil {
			return errors.WithMessage(ErrorGLTFSpec, "NewPlayer.Samplers required")
		}
	}
	return nil
}

func (s *SpecAnimation) To(ctx *parserContext) interface{} {
	res := new(Animation)
	res.Channels = make([]*AnimationChannel, len(s.Channels))
	res.Samplers = make([]*AnimationSampler, len(s.Samplers))
	res.Extensions = s.Extensions
	res.Extras = s.Extras
	return res
}

func (s *SpecAnimation) GetChild(i int) ToGLTF {
	if chleng := len(s.Channels); i < chleng {
		return &s.Channels[i]
	} else {
		return &s.Samplers[i-chleng]
	}
}
func (s *SpecAnimation) SetChild(i int, dst, object interface{}) {
	if chleng := len(s.Channels); i < chleng {
		dst.(*Animation).Channels[i] = object.(*AnimationChannel)
	} else {
		dst.(*Animation).Samplers[i-chleng] = object.(*AnimationSampler)
	}
}
func (s *SpecAnimation) LenChild() int {
	return len(s.Channels) + len(s.Samplers)
}
func (s *SpecAnimation) ImpleGetChild(i int, dst interface{}) interface{} {
	if chleng := len(s.Channels); i < chleng {
		return dst.(*Animation).Channels[i]
	} else {
		return dst.(*Animation).Samplers[i-chleng]
	}
}
