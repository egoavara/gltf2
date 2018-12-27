package gltf2

type ExtensionType interface {
	ExtensionName() string
	Constructor(src []byte) (Specifier, error)
}
//
//func (s *ExtensionKey) String() string {
//	return s.Name
//}
//
//// https://github.com/KhronosGroup/glTF/tree/master/extensions
//var (
//	// Khronos extensions
//	KHR_draco_mesh_compression = &ExtensionKey{
//		Name: "KHR_draco_mesh_compression",
//	}
//	KHR_lights_punctual = &ExtensionKey{
//		Name: "KHR_lights_punctual",
//	}
//
//	KHR_materials_unlit = &ExtensionKey{
//		Name: "KHR_materials_unlit",
//	}
//	KHR_texture_transform = &ExtensionKey{
//		Name: "KHR_texture_transform",
//	}
//	// Vendor extensions
//	ADOBE_materials_thin_transparency = &ExtensionKey{
//		Name: "ADOBE_materials_thin_transparency",
//	}
//	AGI_articulations = &ExtensionKey{
//		Name: "AGI_articulations",
//	}
//	AGI_stk_metadata = &ExtensionKey{
//		Name: "AGI_stk_metadata",
//	}
//	EXT_lights_image_based = &ExtensionKey{
//		Name: "EXT_lights_image_based",
//	}
//	MSFT_lod = &ExtensionKey{
//		Name: "MSFT_lod",
//	}
//	MSFT_texture_dds = &ExtensionKey{
//		Name: "MSFT_texture_dds",
//	}
//	MSFT_packing_normalRoughnessMetallic = &ExtensionKey{
//		Name: "MSFT_packing_normalRoughnessMetallic",
//	}
//	MSFT_packing_occlusionRoughnessMetallic = &ExtensionKey{
//		Name: "MSFT_packing_occlusionRoughnessMetallic",
//	}
//	// In progress
//	// KHR_techniques_webgl
//	// KHR_compressed_texture_transmission
//	// KHR_blend
//	// HDR textures
//	// Advanced PBR materials
//)
//
func extFindByName(name string, exts ...ExtensionType) ExtensionType {
	for _, v := range exts {
		if v.ExtensionName() == name {
			return v
		}
	}
	return nil
}