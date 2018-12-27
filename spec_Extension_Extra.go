package gltf2

type Extensions map[string]ExtensionType

func (s *Extensions) Get(extType ExtensionType) ExtensionType {
	if data, ok := (*s)[extType.ExtensionName()]; ok{
		return data
	}
	return nil
}
func (s *Extensions) GetByName(name string) ExtensionType {
	if data, ok := (*s)[name]; ok{
		return data
	}
	return nil
}
type SpecExtensions map[string]*jsonRawString
type jsonRawString struct {
	src  []byte
	data Specifier
}

func (s *jsonRawString) String() string {
	return string(s.src)
}

func (s *jsonRawString) UnmarshalJSON(src []byte) error {
	s.src = src
	return nil
}

type Extras map[string]interface{}
