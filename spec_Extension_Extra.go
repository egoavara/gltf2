package gltf2

type Extension struct {
	Name string
}

// https://github.com/KhronosGroup/glTF/tree/master/extensions
var (
	// Khronos extensions
	KHR_draco_mesh_compression = &Extension{
		Name: "KHR_draco_mesh_compression",
	}
	KHR_lights_punctual = &Extension{
		Name: "KHR_lights_punctual",
	}
	KHR_materials_pbrSpecularGlossiness = &Extension{
		Name: "KHR_materials_pbrSpecularGlossiness",
	}
	KHR_materials_unlit = &Extension{
		Name: "KHR_materials_unlit",
	}
	KHR_texture_transform = &Extension{
		Name: "KHR_texture_transform",
	}
	// Vendor extensions
	ADOBE_materials_thin_transparency = &Extension{
		Name: "ADOBE_materials_thin_transparency",
	}
	AGI_articulations = &Extension{
		Name: "AGI_articulations",
	}
	AGI_stk_metadata = &Extension{
		Name: "AGI_stk_metadata",
	}
	EXT_lights_image_based = &Extension{
		Name: "EXT_lights_image_based",
	}
	MSFT_lod = &Extension{
		Name: "MSFT_lod",
	}
	MSFT_texture_dds = &Extension{
		Name: "MSFT_texture_dds",
	}
	MSFT_packing_normalRoughnessMetallic = &Extension{
		Name: "MSFT_packing_normalRoughnessMetallic",
	}
	MSFT_packing_occlusionRoughnessMetallic = &Extension{
		Name: "MSFT_packing_occlusionRoughnessMetallic",
	}
	// In progress
	// KHR_techniques_webgl
	// KHR_compressed_texture_transmission
	// KHR_blend
	// HDR textures
	// Advanced PBR materials
)
var (
	exts = []*Extension{
		// TODO KHR_draco_mesh_compression,
		// TODO KHR_lights_punctual,
		KHR_materials_pbrSpecularGlossiness,
		// TODO KHR_materials_unlit,
		// TODO KHR_texture_transform,
		// TODO ADOBE_materials_thin_transparency,
		// TODO AGI_articulations,
		// TODO AGI_stk_metadata,
		// TODO EXT_lights_image_based,
		// TODO MSFT_lod,
		// TODO MSFT_texture_dds,
		// TODO MSFT_packing_normalRoughnessMetallic,
		// TODO MSFT_packing_occlusionRoughnessMetallic,
	}
)

func extName(name string) *Extension {
	for _, v := range exts {
		if v.Name == name {
			return v
		}
	}
	return nil
}
func extNames(name ...string) []*Extension {
	res := make([]*Extension, 0, len(name))
	for _, v := range name {
		res = append(res, extName(v))
	}
	return res
}
func extExist(name string) bool {

	for _, v := range exts {
		if v.Name == name {
			return true
		}
	}
	return false
}
func extExists(name ...string) (bool, string) {
	for _, v := range name {
		if !extExist(v) {
			return false, v
		}
	}
	return true, ""
}

type Extensions map[*Extension]interface{}
type SpecExtensions map[string]JSONRawString
type JSONRawString []byte

func (s JSONRawString) String() string {
	return string(s)
}

func (s *JSONRawString) UnmarshalJSON(src []byte) error {
	*s = src
	return nil
}

//func (s *SpecExtensions) UnmarshalJSON(src []byte) error {
//	s.src = src
//	if src != nil{
//		fmt.Println(string(src))
//	}
//	return nil
//}

type Extras map[string]interface{}
