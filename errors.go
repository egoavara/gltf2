package gltf2

import "github.com/pkg/errors"

var (
	ErrorJSON         = errors.New("Json parsing fail")
	ErrorEnum         = errors.New("Enum parsing fail")
	ErrorGLTFSpec     = errors.New("Specifier fail")
	ErrorGLTFLink     = errors.New("glTFid link not found")
	ErrorTask         = errors.New("Task fail")
	ErrorParser       = errors.New("Parser error")
	ErrorParserOption = errors.New("Parser option error")
)
