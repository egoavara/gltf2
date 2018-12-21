package gltf2

type ToGLTF interface {
	Scheme() string
	Syntax(strictness Strictness, root interface{}) error
	To(ctx *parserContext) interface{}
}
type HaveExtensions interface {
	GetExtensions() *Extensions

}
type ChildrunToGLTF interface {
	ToGLTF

	GetChild(i int) ToGLTF
	SetChild(i int, dst, object interface{})
	LenChild() int
	ImpleGetChild(i int, dst interface{}) interface{}
}
type LinkToGLTF interface {
	ToGLTF
	// s.Link(Root, s.To())
	//              ^ important!
	Link(Root *GLTF, parent interface{}, dst interface{}) error
}

const (
	SCHEME_GLTF                            = "glTF"
	SCHEME_BUFFER                          = "buffer"
	SCHEME_BUFFERVIEW                      = "bufferView"
	SCHEME_ACCESSOR                        = "accessor"
	SCHEME_ASSET                           = "asset"
	SCHEME_CAMERA                          = "camera"
	SCHEME_CAMERA_PERSPECTIVE              = "camera/perspective"
	SCHEME_CAMERA_ORTHOGRAPHIC             = "camera/orthographic"
	SCHEME_IMAGE                           = "image"
	SCHEME_SAMPLER                         = "sampler"
	SCHEME_TEXTURE                         = "texture"
	SCHEME_TEXTURE_INFO                    = "texture/info"
	SCHEME_MATERIAL                        = "material"
	SCHEME_MATERIAL_PBR_METALLIC_ROUGHNESS = "material/pbrMetallicRoughness"
	SCHEME_MATERIAL_OCCLUSION_TEXTUREINFO  = "material/occlusionTextureInfo"
	SCHEME_MATERIAL_NORMAL_TEXTUREINFO     = "material/normalTextureInfo"
	SCHEME_MESH                            = "mesh"
	SCHEME_MESH_PRIMITIVE                  = "mesh/primitive"
	SCHEME_NODE                            = "node"
	SCHEME_SCENE                           = "scene"
	SCHEME_ANIMATION                       = "animation"
	SCHEME_ANIMATION_SAMPLER               = "animation/sampler"
	SCHEME_ANIMATION_CHANNEL               = "animation/channel"
	SCHEME_ANIMATION_CHANNEL_TARGET        = "animation/channel/target"
	//
	SCHEME_EXTENSION = "extension"
)
