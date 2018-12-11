package gltf2

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
)

type URI url.URL

func (s *URI) String() string {
	return s.Data().String()
}

func (s *URI) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	res, err := url.Parse(v)
	if err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	*s = URI(*res)
	return nil
}

func (s *URI) Data() *url.URL {
	res := new(url.URL)
	*res = url.URL(*s)
	return res
}
