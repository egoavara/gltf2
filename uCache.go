package gltf2

func ThrowAllCache(gltf *GLTF)  {
	for _, b := range gltf.Buffers{
		b.ThrowCache()
	}
	for _, i := range gltf.Images{
		i.ThrowCache()
	}
}
