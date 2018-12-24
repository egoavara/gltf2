package gltf2

type Specifier interface {
	Scheme() string
	Syntax(strictness Strictness, root interface{}) error
	To(ctx *parserContext) interface{}

}
type ExtensionGetter interface {
	GetExtension() *SpecExtensions
}
type ExtensionSetter interface {
	SetExtension(extensions *Extensions)
}
type Parents interface {
	Specifier

	GetChild(i int) Specifier
	SetChild(i int, dst, object interface{})
	LenChild() int
	ImpleGetChild(i int, dst interface{}) interface{}
}
type Linker interface {
	Specifier
	// s.Link(Root, s.To())
	//              ^ important!
	Link(Root *GLTF, parent interface{}, dst interface{}) error
}
