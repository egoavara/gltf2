package gltf2

import (
	"github.com/pkg/errors"
)

type AnimationChannelTarget struct {
	Node       *Node
	Path       Path
	Extensions *Extensions
	Extras     *Extras
}

type SpecAnimationChannelTarget struct {
	Node       *SpecGLTFID `json:"node"`
	Path       *Path       `json:"path"` // required
	Extensions *Extensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`
}

func (s *SpecAnimationChannelTarget) GetExtensions() *Extensions {
	return s.Extensions
}

func (s *SpecAnimationChannelTarget) Scheme() string {
	return SCHEME_ANIMATION_CHANNEL_TARGET
}
func (s *SpecAnimationChannelTarget) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		// TODO
		// https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#animations target.Path constraint
		fallthrough
	case LEVEL1:
		if s.Path == nil {
			return errors.Errorf("AnimationChannelTarget.Path required")
		}
	}
	return nil
}
func (s *SpecAnimationChannelTarget) To(ctx *parserContext) interface{} {
	res := new(AnimationChannelTarget)
	res.Path = *s.Path
	res.Extensions = s.Extensions
	res.Extras = s.Extras
	return res
}
func (s *SpecAnimationChannelTarget) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if s.Node != nil {
		if !inRange(*s.Node, len(Root.Nodes)) {
			return errors.Errorf("AnimationChannelTarget.Node linking fail")
		}
		dst.(*AnimationChannelTarget).Node = Root.Nodes[*s.Node]
	}
	return nil
}
