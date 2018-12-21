package gltf2

func isUniqueExtension(datas ...string) (bool, string) {
	for i := 0; i < len(datas); i++ {
		for j := i + 1; j < len(datas); j++ {
			if datas[i] == datas[j] {
				return false, datas[i]
			}
		}
	}
	return true, ""
}

func isUniqueGLTFID(datas ...SpecGLTFID) (bool, SpecGLTFID) {
	for i := 0; i < len(datas); i++ {
		for j := i + 1; j < len(datas); j++ {
			if datas[i] == datas[j] {
				return false, datas[i]
			}
		}
	}
	return true, 0
}
