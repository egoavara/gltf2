package gltf2

type Specifier interface {
	Scheme() string
	Syntax(strictness Strictness, root Specifier, parent Specifier) error
	To(ctx *parserContext) interface{}
}
type ExtensionSpecifier interface {
	SpecExtension() *SpecExtensions
}
type ExtensionStructure interface {
	SetExtension(extensions *Extensions)
	GetExtension() *Extensions
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
