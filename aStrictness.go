package gltf2

type Strictness uint8

// Anything other than LEVEL{0 to 3} is treated as LEVLE0
const (
	// Do not check anything
	//
	LEVEL0 Strictness = 0
	// Check essential only
	//
	// + required field
	// + dependency
	LEVEL1 Strictness = 1
	// Check potential problematic
	//
	// + unique slice item
	// + slice item count
	// ! Camera/Perspective : + Argument large, etc
	// ! Camera/Orthographic : + Argument large, etc
	// ! Material/pbrMetallicRoughness : + Argument range, etc
	LEVEL2 Strictness = 2
	// Follow Specification
	//
	//
	// + max, min limitation
	LEVEL3 Strictness = 3
)
