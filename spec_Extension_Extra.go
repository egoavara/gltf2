package gltf2

type Extension string

// https://github.com/KhronosGroup/glTF/tree/master/extensions
const (
	// Khronos extensions
	KHR_draco_mesh_compression          Extension = "KHR_draco_mesh_compression"
	KHR_lights_punctual                 Extension = "KHR_lights_punctual"
	KHR_materials_pbrSpecularGlossiness Extension = "KHR_materials_pbrSpecularGlossiness"
	KHR_materials_unlit                 Extension = "KHR_materials_unlit"
	KHR_texture_transform               Extension = "KHR_texture_transform"
	// Vendor extensions
	ADOBE_materials_thin_transparency       Extension = "ADOBE_materials_thin_transparency"
	AGI_articulations                       Extension = "AGI_articulations"
	AGI_stk_metadata                        Extension = "AGI_stk_metadata"
	EXT_lights_image_based                  Extension = "EXT_lights_image_based"
	MSFT_lod                                Extension = "MSFT_lod"
	MSFT_texture_dds                        Extension = "MSFT_texture_dds"
	MSFT_packing_normalRoughnessMetallic    Extension = "MSFT_packing_normalRoughnessMetallic"
	MSFT_packing_occlusionRoughnessMetallic Extension = "MSFT_packing_occlusionRoughnessMetallic"
	// In progress
	// KHR_techniques_webgl
	// KHR_compressed_texture_transmission
	// KHR_blend
	// HDR textures
	// Advanced PBR materials
)

type Extensions map[string]interface{}
type Extras map[string]interface{}
