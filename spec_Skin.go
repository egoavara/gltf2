package gltf2

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Skin struct {
	InverseBindMatrices *SpecGLTFID  `json:"inverseBindMatrices"`
	Skeleton            *SpecGLTFID  `json:"skeleton"`
	Joints              []SpecGLTFID `json:"joints"` // require, minItem(1), unique
	Name                *string      `json:"name,omitempty"`
	Extensions          *Extensions  `json:"extensions,omitempty"`
	Extras              *Extras      `json:"extras,omitempty"`
}

func (s *Skin) UnmarshalJSON(data []byte) error {
	type Temp struct {
		InverseBindMatrices *SpecGLTFID  `json:"inverseBindMatrices"`
		Skeleton            *SpecGLTFID  `json:"skeleton"`
		Joints              []SpecGLTFID `json:"joints"` // require, unique
		Name                *string      `json:"name,omitempty"`
		Extensions          *Extensions  `json:"extensions,omitempty"`
		Extras              *Extras      `json:"extras,omitempty"`
	}

	var temp = Temp{}
	// Parse
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.WithMessage(ErrorJSON, err.Error())
	}
	// Spec
	if temp.Joints == nil {
		return errors.WithMessage(ErrorGLTFSpec, "Skin.Joints require + minItem(1)")
	}
	//if !isUniqueGLTFID(temp.Joints...) {
	//	return errors.WithMessage(ErrorGLTFSpec, "Skin.Joints unique")
	//}
	// _Setup
	s.InverseBindMatrices = temp.InverseBindMatrices
	s.Skeleton = temp.Skeleton
	s.Joints = temp.Joints
	s.Name = temp.Name
	s.Extensions = temp.Extensions
	s.Extras = temp.Extras
	return nil
}
