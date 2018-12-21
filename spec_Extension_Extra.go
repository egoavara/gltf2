package gltf2

type Extension struct {
	Name          string
	Constructor   func() interface{}
	Specification func() interface{}
}

// https://github.com/KhronosGroup/glTF/tree/master/extensions
var (
	// Khronos extensions
	KHR_draco_mesh_compression = &Extension{
		Name: "KHR_draco_mesh_compression",
		Constructor: func() interface{} {
			return nil
		},
	}
	KHR_lights_punctual = &Extension{
		Name: "KHR_lights_punctual",
		Constructor: func() interface{} {
			return nil
		},
	}
	KHR_materials_pbrSpecularGlossiness = &Extension{
		Name: "KHR_materials_pbrSpecularGlossiness",
		Constructor: func() interface{} {
			return nil
		},
		Specification: func() interface{} {
			return nil
		},
	}
	KHR_materials_unlit = &Extension{
		Name: "KHR_materials_unlit",
		Constructor: func() interface{} {
			return nil
		},
	}
	KHR_texture_transform = &Extension{
		Name: "KHR_texture_transform",
		Constructor: func() interface{} {
			return nil
		},
	}
	// Vendor extensions
	ADOBE_materials_thin_transparency = &Extension{
		Name: "ADOBE_materials_thin_transparency",
		Constructor: func() interface{} {
			return nil
		},
	}
	AGI_articulations = &Extension{
		Name: "AGI_articulations",
		Constructor: func() interface{} {
			return nil
		},
	}
	AGI_stk_metadata = &Extension{
		Name: "AGI_stk_metadata",
		Constructor: func() interface{} {
			return nil
		},
	}
	EXT_lights_image_based = &Extension{
		Name: "EXT_lights_image_based",
		Constructor: func() interface{} {
			return nil
		},
	}
	MSFT_lod = &Extension{
		Name: "MSFT_lod",
		Constructor: func() interface{} {
			return nil
		},
	}
	MSFT_texture_dds = &Extension{
		Name: "MSFT_texture_dds",
		Constructor: func() interface{} {
			return nil
		},
	}
	MSFT_packing_normalRoughnessMetallic = &Extension{
		Name: "MSFT_packing_normalRoughnessMetallic",
		Constructor: func() interface{} {
			return nil
		},
	}
	MSFT_packing_occlusionRoughnessMetallic = &Extension{
		Name: "MSFT_packing_occlusionRoughnessMetallic",
		Constructor: func() interface{} {
			return nil
		},
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
type Extras map[*Extension]interface{}
