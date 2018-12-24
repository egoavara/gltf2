package gltf2

import (
	"github.com/pkg/errors"
)

// https://github.com/KhronosGroup/glTF/blob/master/specification/2.0/schema/glTF.schema.json
type GLTF struct {
	ExtensionsUsed     []*Extension
	ExtensionsRequired []*Extension
	Accessors          []*Accessor
	Asset              Asset
	Buffers            []*Buffer
	BufferViews        []*BufferView
	Cameras            []*Camera
	Images             []Image
	Materials          []*Material
	Meshes             []*Mesh
	Nodes              []*Node
	Samplers           []*Sampler
	Scene              *Scene
	Scenes             []*Scene
	Textures           []*Texture
	Animations         []*Animation
	Extensions         *Extensions
	Extras             *Extras

	// None spec
	UserData interface{}
}

func (s *GLTF) SetExtension(extensions *Extensions) {
	s.Extensions = extensions
}

type SpecGLTF struct {
	ExtensionsUsed     []string         `json:"extensionsUsed,omitempty"`     // minitem(1), unique
	ExtensionsRequired []string         `json:"extensionsRequired,omitempty"` // minitem(1), unique
	Accessors          []SpecAccessor   `json:"accessors,omitempty"`          // minitem(1)
	Asset              *SpecAsset       `json:"asset,omitempty"`              // required
	Buffers            []SpecBuffer     `json:"buffers,omitempty"`            // minitem(1)
	BufferViews        []SpecBufferView `json:"bufferViews,omitempty"`        // minitem(1)
	Cameras            []SpecCamera     `json:"cameras,omitempty"`            // minitem(1)
	Images             []SpecImage      `json:"images,omitempty"`             // minitem(1)
	Materials          []SpecMaterial   `json:"materials,omitempty"`          // minitem(1)
	Meshes             []SpecMesh       `json:"meshes,omitempty"`             // minitem(1)
	Nodes              []SpecNode       `json:"nodes,omitempty"`              // minitem(1)
	Samplers           []SpecSampler    `json:"samplers,omitempty"`           // minitem(1)
	Scene              *SpecGLTFID      `json:"scene,omitempty"`              // dependency(scenes)
	Scenes             []SpecScene      `json:"scenes,omitempty"`             // minitem(1)
	Textures           []SpecTexture    `json:"textures,omitempty"`           // minitem(1)
	Animations         []SpecAnimation  `json:"animations,omitempty"`         // minitem(1)
	// TODO : Skins              []Skin       `json:"skins,omitempty"`              // minitem(1)
	Extensions *SpecExtensions `json:"extensions,omitempty"`
	Extras     *Extras         `json:"extras,omitempty"`
	//
	cache []int
}

func (s *SpecGLTF) buildSchemeIndexCache() {
	var lengths []int
	var temp = 0
	lengths = append(lengths, temp)

	temp += len(s.Accessors)
	lengths = append(lengths, temp)

	temp += 1 // s.Asset
	lengths = append(lengths, temp)

	temp += len(s.Buffers)
	lengths = append(lengths, temp)

	temp += len(s.BufferViews)
	lengths = append(lengths, temp)

	temp += len(s.Cameras)
	lengths = append(lengths, temp)

	temp += len(s.Images)
	lengths = append(lengths, temp)

	temp += len(s.Materials)
	lengths = append(lengths, temp)

	temp += len(s.Meshes)
	lengths = append(lengths, temp)

	temp += len(s.Nodes)
	lengths = append(lengths, temp)

	temp += len(s.Samplers)
	lengths = append(lengths, temp)

	temp += len(s.Scenes)
	lengths = append(lengths, temp)

	temp += len(s.Textures)
	lengths = append(lengths, temp)

	temp += len(s.Animations)
	lengths = append(lengths, temp)
	//
	s.cache = lengths
}
func (s *SpecGLTF) getSchemeIndex(i int) (scheme string, schemeIndex int) {
	if s.cache == nil {
		s.buildSchemeIndexCache()
	}
	for cacheIndex := 0; cacheIndex < len(s.cache)-1; cacheIndex++ {
		if i < s.cache[cacheIndex+1] {
			schemeIndex = i - s.cache[cacheIndex]
			switch cacheIndex {
			case 0:
				scheme = SCHEME_ACCESSOR
			case 1:
				scheme = SCHEME_ASSET
			case 2:
				scheme = SCHEME_BUFFER
			case 3:
				scheme = SCHEME_BUFFERVIEW
			case 4:
				scheme = SCHEME_CAMERA
			case 5:
				scheme = SCHEME_IMAGE
			case 6:
				scheme = SCHEME_MATERIAL
			case 7:
				scheme = SCHEME_MESH
			case 8:
				scheme = SCHEME_NODE
			case 9:
				scheme = SCHEME_SAMPLER
			case 10:
				scheme = SCHEME_SCENE
			case 11:
				scheme = SCHEME_TEXTURE
			case 12:
				scheme = SCHEME_ANIMATION
			}
			break
		}
	}
	return
}

func (s *SpecGLTF) GetExtension() *SpecExtensions {
	return s.Extensions
}
func (s *SpecGLTF) Scheme() string {
	return SCHEME_GLTF
}
func (s *SpecGLTF) Syntax(strictness Strictness, root interface{}) error {
	switch strictness {
	case LEVEL3:
		fallthrough
	case LEVEL2:
		if ok, data := isUniqueExtension(s.ExtensionsUsed...); !ok {
			return errors.Errorf("GLTF.ExtensionsUsed is unique, but duplicate item '%s'", data)
		}
		if ok, data := isUniqueExtension(s.ExtensionsRequired...); !ok {
			return errors.Errorf("GLTF.ExtensionsUsed is unique, but duplicate item '%s'", data)
		}
		fallthrough
	case LEVEL1:
		if s.Scene != nil {
			if s.Scenes == nil {
				return errors.WithMessage(ErrorGLTFSpec, "GLTF.Scene dependency(scenes)")
			}
		}
		if s.Asset == nil {
			return errors.Errorf("GLTF.Asset required")
		}
		if ok, name := extExists(s.ExtensionsRequired...); !ok {
			return errors.Errorf("GLTF.ExtensionsRequired has unknown '%s'", name)
		}
	}
	return nil
}
func (s *SpecGLTF) To(ctx *parserContext) interface{} {
	res := new(GLTF)
	res.ExtensionsUsed = extNames(s.ExtensionsUsed...)
	res.ExtensionsRequired = extNames(s.ExtensionsUsed...)
	//
	res.Accessors = make([]*Accessor, len(s.Accessors))
	res.Buffers = make([]*Buffer, len(s.Buffers))
	res.BufferViews = make([]*BufferView, len(s.BufferViews))
	res.Cameras = make([]*Camera, len(s.Cameras))
	res.Images = make([]Image, len(s.Images))
	res.Materials = make([]*Material, len(s.Materials))
	res.Meshes = make([]*Mesh, len(s.Meshes))
	res.Nodes = make([]*Node, len(s.Nodes))
	res.Samplers = make([]*Sampler, len(s.Samplers))
	res.Scenes = make([]*Scene, len(s.Scenes))
	res.Textures = make([]*Texture, len(s.Textures))
	res.Animations = make([]*Animation, len(s.Animations))
	//
	res.Extras = s.Extras
	return res
}
func (s *SpecGLTF) Link(Root *GLTF, parent interface{}, dst interface{}) error {
	if s.Scene != nil {
		if !inRange(*s.Scene, len(Root.Scenes)) {
			return errors.Errorf("glTF.Scene linking fail")
		}
		dst.(*GLTF).Scene = Root.Scenes[*s.Scene]
	}
	return nil
}
func (s *SpecGLTF) GetChild(i int) Specifier {
	scheme, schemeIndex := s.getSchemeIndex(i)
	switch scheme {
	case SCHEME_ACCESSOR:
		return &s.Accessors[schemeIndex]
	case SCHEME_ASSET:
		return s.Asset
	case SCHEME_BUFFER:
		return &s.Buffers[schemeIndex]
	case SCHEME_BUFFERVIEW:
		return &s.BufferViews[schemeIndex]
	case SCHEME_CAMERA:
		return &s.Cameras[schemeIndex]
	case SCHEME_IMAGE:
		return &s.Images[schemeIndex]
	case SCHEME_MATERIAL:
		return &s.Materials[schemeIndex]
	case SCHEME_MESH:
		return &s.Meshes[schemeIndex]
	case SCHEME_NODE:
		return &s.Nodes[schemeIndex]
	case SCHEME_SAMPLER:
		return &s.Samplers[schemeIndex]
	case SCHEME_SCENE:
		return &s.Scenes[schemeIndex]
	case SCHEME_TEXTURE:
		return &s.Textures[schemeIndex]
	case SCHEME_ANIMATION:
		return &s.Animations[schemeIndex]
	}
	return nil
}
func (s *SpecGLTF) SetChild(i int, dst, object interface{}) {
	scheme, schemeIndex := s.getSchemeIndex(i)
	switch scheme {
	case SCHEME_ACCESSOR:
		dst.(*GLTF).Accessors[schemeIndex] = object.(*Accessor)
	case SCHEME_ASSET:
		dst.(*GLTF).Asset = *object.(*Asset)
	case SCHEME_BUFFER:
		dst.(*GLTF).Buffers[schemeIndex] = object.(*Buffer)
	case SCHEME_BUFFERVIEW:
		dst.(*GLTF).BufferViews[schemeIndex] = object.(*BufferView)
	case SCHEME_CAMERA:
		dst.(*GLTF).Cameras[schemeIndex] = object.(*Camera)
	case SCHEME_IMAGE:
		dst.(*GLTF).Images[schemeIndex] = object.(Image)
	case SCHEME_MATERIAL:
		dst.(*GLTF).Materials[schemeIndex] = object.(*Material)
	case SCHEME_MESH:
		dst.(*GLTF).Meshes[schemeIndex] = object.(*Mesh)
	case SCHEME_NODE:
		dst.(*GLTF).Nodes[schemeIndex] = object.(*Node)
	case SCHEME_SAMPLER:
		dst.(*GLTF).Samplers[schemeIndex] = object.(*Sampler)
	case SCHEME_SCENE:
		dst.(*GLTF).Scenes[schemeIndex] = object.(*Scene)
	case SCHEME_TEXTURE:
		dst.(*GLTF).Textures[schemeIndex] = object.(*Texture)
	case SCHEME_ANIMATION:
		dst.(*GLTF).Animations[schemeIndex] = object.(*Animation)
	}
}
func (s *SpecGLTF) LenChild() int {
	if s.cache == nil {
		s.buildSchemeIndexCache()
	}
	return s.cache[len(s.cache)-1]
}
func (s *SpecGLTF) ImpleGetChild(i int, dst interface{}) interface{} {
	scheme, schemeIndex := s.getSchemeIndex(i)
	switch scheme {
	case SCHEME_ACCESSOR:
		return dst.(*GLTF).Accessors[schemeIndex]
	case SCHEME_ASSET:
		return dst.(*GLTF).Asset
	case SCHEME_BUFFER:
		return dst.(*GLTF).Buffers[schemeIndex]
	case SCHEME_BUFFERVIEW:
		return dst.(*GLTF).BufferViews[schemeIndex]
	case SCHEME_CAMERA:
		return dst.(*GLTF).Cameras[schemeIndex]
	case SCHEME_IMAGE:
		return dst.(*GLTF).Images[schemeIndex]
	case SCHEME_MATERIAL:
		return dst.(*GLTF).Materials[schemeIndex]
	case SCHEME_MESH:
		return dst.(*GLTF).Meshes[schemeIndex]
	case SCHEME_NODE:
		return dst.(*GLTF).Nodes[schemeIndex]
	case SCHEME_SAMPLER:
		return dst.(*GLTF).Samplers[schemeIndex]
	case SCHEME_SCENE:
		return dst.(*GLTF).Scenes[schemeIndex]
	case SCHEME_TEXTURE:
		return dst.(*GLTF).Textures[schemeIndex]
	case SCHEME_ANIMATION:
		return dst.(*GLTF).Animations[schemeIndex]
	}
	return nil
}
