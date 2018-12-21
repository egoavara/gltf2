package gltf2

import (
	"encoding/json"
	"fmt"
	"github.com/iamGreedy/essence/version"

	"github.com/pkg/errors"
)

type Asset struct {
	Copyright  string
	Generator  string
	Version    version.Version
	MinVersion *version.Version
	Extensions *Extensions
	Extras     *Extras
}

func (s *Asset) String() string {
	bts, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("%e", err)
	}
	return string(bts)
}

type SpecAsset struct {
	Copyright  *string          `json:"copyright,omitempty"`
	Generator  *string          `json:"generator,omitempty"`
	Version    *version.Version `json:"version"` // required
	MinVersion *version.Version `json:"minVersion,omitempty"`
	Extensions *Extensions      `json:"extensions,omitempty"`
	Extras     *Extras          `json:"extras,omitempty"`
}

func (s *SpecAsset) GetExtensions() *Extensions {
	return s.Extensions
}

func (s *SpecAsset) Scheme() string {
	return SCHEME_ASSET
}
func (s *SpecAsset) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		fallthrough
	case LEVEL1:
		if s.Version == nil {
			return errors.Errorf("Asset.Version required")
		}
	}
	return nil
}
func (s *SpecAsset) To(ctx *parserContext) interface{} {
	res := new(Asset)
	if s.Copyright != nil {
		res.Copyright = *s.Copyright
	}
	if s.Generator != nil {
		res.Generator = *s.Generator
	}
	res.Version = *s.Version
	res.MinVersion = s.MinVersion
	res.Extras = s.Extras
	res.Extensions = s.Extensions
	return res
}
