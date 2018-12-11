package gltf2

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type AccessorSparse struct {
	Count   int                   `json:"count"`  // required, minimum(0)
	Indices AccessorSparseIndices `json:"indice"` // required
	Values  AccessorSparseValues  `json:"values"` // required
}

func (s *AccessorSparse) UnmarshalJSON(data []byte) error {
	type Temp struct {
		Count   *int                   `json:"count"`  // required, minimum(0)
		Indices *AccessorSparseIndices `json:"indice"` // required
		Values  *AccessorSparseValues  `json:"values"` // required
	}
	var temp Temp
	// Parse
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	// Spec
	if temp.Count == nil {
		return errors.WithMessage(ErrorGLTFSpec, "AccessorSparse.Count required")
	}
	if temp.Indices == nil {
		return errors.WithMessage(ErrorGLTFSpec, "AccessorSparse.Indices required")
	}
	if temp.Values == nil {
		return errors.WithMessage(ErrorGLTFSpec, "AccessorSparse.Values required")
	}
	//_Setup
	s.Count = *temp.Count
	s.Indices = *temp.Indices
	s.Values = *temp.Values
	return nil
}
