package gltf2

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/glTFid.schema.json
type SpecGLTFID int

func (s *SpecGLTFID) UnmarshalJSON(src []byte) error {
	var temp int
	if err := json.Unmarshal(src, &temp); err != nil{return err}
	if temp < 0{
		return errors.Errorf("glTFid minimum is 0, but %d", *s)
	}
	*s = SpecGLTFID(temp)
	return nil
}

func (s SpecGLTFID) String() string {
	return fmt.Sprintf("glTFid(%d)", s)
}
