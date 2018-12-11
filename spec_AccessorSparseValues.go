package gltf2

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type AccessorSparseValues struct {
	BufferView SpecGLTFID  `json:"bufferView"`           // required
	ByteOffset *int        `json:"byteOffset,omitempty"` // default(0), minimum(0)
	Extensions *Extensions `json:"extensions,omitempty"`
	Extras     *Extras     `json:"extras,omitempty"`
}

func (s *AccessorSparseValues) UnmarshalJSON(data []byte) error {
	type Temp struct {
		BufferView *SpecGLTFID `json:"bufferView"`           // required
		ByteOffset *int        `json:"byteOffset,omitempty"` // default(0), minimum(0)
		Extensions *Extensions `json:"extensions,omitempty"`
		Extras     *Extras     `json:"extras,omitempty"`
	}
	var temp Temp
	// Parse
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	// Spec
	if temp.BufferView == nil {
		return errors.WithMessage(ErrorGLTFSpec, "AccessorSparseIndices.BufferView required")
	}
	if temp.ByteOffset != nil {
		if *temp.ByteOffset < 0 {
			return errors.WithMessage(ErrorGLTFSpec, "AccessorSparseIndices.ByteOffset minimum(0)")
		}
	}
	// _Setup
	s.BufferView = *temp.BufferView
	s.ByteOffset = temp.ByteOffset
	s.Extensions = temp.Extensions
	s.Extras = temp.Extras
	return nil
}
