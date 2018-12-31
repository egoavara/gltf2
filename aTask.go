package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/align"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/glog"
	"github.com/labstack/gommon/bytes"
	"github.com/pkg/errors"
	"math"
	"net/url"
)

type (
	Task interface {
		TaskName() string
	}
	PreTask interface {
		Task
		PreLoad(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger)
	}
	PostTask interface {
		Task
		PostLoad(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error
	}
)

type (
	_fn_preonly struct {
		name string
		fn   func(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger)
	}
	_fn_postonly struct {
		name string
		fn   func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error
	}
)

func (s *_fn_preonly) TaskName() string {
	return s.name
}
func (s *_fn_preonly) PreLoad(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger) {
	s.fn(parser, gltf, logger)
}

func (s *_fn_postonly) TaskName() string {
	return s.name
}
func (s *_fn_postonly) PostLoad(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
	return s.fn(parser, gltf, logger)
}

func FnPreTask(name string, fn func(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger)) Task {
	return &_fn_preonly{
		name: name,
		fn:   fn,
	}
}
func FnPostTask(name string, fn func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error) Task {
	return &_fn_postonly{
		name: name,
		fn:   fn,
	}
}

func _BufferCaching(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
	for i, buf := range gltf.Buffers {
		if len(buf.Name) > 0 {
			logger.Printf("glTF.Buffers['%s'] ", buf.Name)
		} else {
			logger.Printf("glTF.Buffers[%d] ", i)
		}
		//
		inner := logger.Indent()
		if buf.IsCached() {
			inner.Println("Already cached")
			continue
		}
		uesc, _ := url.PathUnescape(buf.URI.Data().String())
		if buf.ByteLength == nil {
			inner.Printf("Caching : %s\n", uesc)
		} else {
			inner.Printf("Caching : %s, total = %s\n", uesc, bytes.Format(int64(*buf.ByteLength)))
		}
		if _, err := buf.Load(true); err != nil {
			return err
		}
		inner.Printf("Cached : total = %s\n", bytes.Format(int64(len(buf.Cache()))))
	}
	return nil
}
func _ImageCaching(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
	for i, img := range gltf.Images {
		if len(img.Name()) > 0 {
			logger.Printf("glTF.Images['%s'] ", img.Name())
		} else {
			logger.Printf("glTF.Images[%d] ", i)
		}
		//
		inner := logger.Indent()
		if img.IsCached() {
			inner.Printf("Already cached\n")
			continue
		}
		switch base := img.(type) {
		case *BufferImage:
			inner.Printf("Caching : BufferView(%v)", base.BufferView)
		case *URIImage:
			uesc, _ := url.PathUnescape(base.URI.Data().String())
			inner.Printf("Caching : %s", uesc)
		}
		if _, err := img.Load(true); err != nil {
			return err
		}
		inner.Printf("Cached : image(rect = %v, size = %s)\n", img.Cache().Rect, bytes.Format(int64(len(img.Cache().Pix))))
	}
	return nil
}

var Tasks = struct {
	HelloWorld Task
	ByeWorld   Task
	// Buffer upload to memory(RAM)
	BufferCaching Task
	// Image upload to memory(RAM)
	ImageCaching Task
	// BufferCaching + ImageCaching
	// If there is Caching, you don't need to call BufferCaching or ImageCaching
	Caching Task

	AccessorMinMax Task
	ModelAlign     func(x, y, z align.Align) Task
	ModelScale     func(axis axis.Axis, meter meter.Meter) Task
	// Set bufferView.Target NEED_TO_DEFINE_BUFFER to real gl.h enum
	AutoBufferTarget Task
	// make accessor no stride
	TightPacking Task
	// TODO UnpackDracoMesh
	// TODO Clean Task
	// - Dangling Node
	// - Unreferenced Buffer, Image
	// TODO :MergeBuffer Task
	// Merge all buffers if there are many buffers
	// It is separated by Accessor
	// Task<MergeBuffer> exist, Parser occur error
	// TODO :SplitBuffer Task
	// Make all Image store in single buffer
	//
	// TODO : BuildNormal Task
	// TODO : BuildTangent Task
}{
	// Simple Task
	HelloWorld: FnPreTask("Hello, World", func(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger) {
		logger.Println("Hello, World")
	}),
	// Simple Task
	ByeWorld: FnPostTask("Bye, World", func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
		logger.Println("Bye, World")
		return nil
	}),
	// Caching Task
	BufferCaching: FnPostTask("Buffer Caching", _BufferCaching),
	// Caching Task
	ImageCaching: FnPostTask("Image Caching", _ImageCaching),
	// Caching Task
	Caching: FnPostTask("Caching", func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
		if err := _BufferCaching(parser, gltf, logger); err != nil {
			return err
		}
		if err := _ImageCaching(parser, gltf, logger); err != nil {
			return err
		}
		return nil
	}),
	// Accessor Task
	AccessorMinMax: FnPostTask("Accessor Min Max", func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
		for i, accessor := range gltf.Accessors {
			if len(accessor.Min) > 0 && len(accessor.Max) > 0 {
				continue
			}
			if accessor.BufferView == nil {
				continue
			}
			var (
				min = make([]float32, accessor.Type.Count())
				max = make([]float32, accessor.Type.Count())
			)
			switch accessor.ComponentType {
			case BYTE:
				var (
					tempmin = make([]int8, accessor.Type.Count())
					tempmax = make([]int8, accessor.Type.Count())
				)
				slice := accessor.MustSliceMapping(new([]int8), false, true).([]int8)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < tempmin[j] {
							tempmin[j] = slice[i+j]
						}
						if slice[i+j] > tempmax[j] {
							tempmax[j] = slice[i+j]
						}
					}
				}
				for i, v := range tempmin {
					min[i] = float32(v)
				}
				for i, v := range tempmax {
					max[i] = float32(v)
				}
			case UNSIGNED_BYTE:
				var (
					tempmin = make([]uint8, accessor.Type.Count())
					tempmax = make([]uint8, accessor.Type.Count())
				)
				slice := accessor.MustSliceMapping(new([]uint8), false, true).([]uint8)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < tempmin[j] {
							tempmin[j] = slice[i+j]
						}
						if slice[i+j] > tempmax[j] {
							tempmax[j] = slice[i+j]
						}
					}
				}
				for i, v := range tempmin {
					min[i] = float32(v)
				}
				for i, v := range tempmax {
					max[i] = float32(v)
				}
			case SHORT:
				var (
					tempmin = make([]int16, accessor.Type.Count())
					tempmax = make([]int16, accessor.Type.Count())
				)
				slice := accessor.MustSliceMapping(new([]int16), false, true).([]int16)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < tempmin[j] {
							tempmin[j] = slice[i+j]
						}
						if slice[i+j] > tempmax[j] {
							tempmax[j] = slice[i+j]
						}
					}
				}
				for i, v := range tempmin {
					min[i] = float32(v)
				}
				for i, v := range tempmax {
					max[i] = float32(v)
				}
			case UNSIGNED_SHORT:
				var (
					tempmin = make([]uint16, accessor.Type.Count())
					tempmax = make([]uint16, accessor.Type.Count())
				)
				slice := accessor.MustSliceMapping(new([]uint16), false, true).([]uint16)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < tempmin[j] {
							tempmin[j] = slice[i+j]
						}
						if slice[i+j] > tempmax[j] {
							tempmax[j] = slice[i+j]
						}
					}
				}
				for i, v := range tempmin {
					min[i] = float32(v)
				}
				for i, v := range tempmax {
					max[i] = float32(v)
				}
			case UNSIGNED_INT:
				var (
					tempmin = make([]uint32, accessor.Type.Count())
					tempmax = make([]uint32, accessor.Type.Count())
				)
				slice := accessor.MustSliceMapping(new([]uint32), false, true).([]uint32)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < tempmin[j] {
							tempmin[j] = slice[i+j]
						}
						if slice[i+j] > tempmax[j] {
							tempmax[j] = slice[i+j]
						}
					}
				}
				for i, v := range tempmin {
					min[i] = float32(v)
				}
				for i, v := range tempmax {
					max[i] = float32(v)
				}
			case FLOAT:
				slice := accessor.MustSliceMapping(new([]float32), false, true).([]float32)
				for i := 0; i < len(slice); i += accessor.Type.Count() {
					for j := 0; j < accessor.Type.Count(); j++ {
						if slice[i+j] < min[j] {
							min[j] = slice[i+j]
						}
						if slice[i+j] > max[j] {
							max[j] = slice[i+j]
						}
					}
				}
			default:
				logger.Panic("Unreachable")
			}

			if len(accessor.Name) > 0 {
				logger.Printf("glTF.Accessors['%s'] ", accessor.Name)
			} else {
				logger.Printf("glTF.Accessors[%d] ", i)
			}
			//
			inner := logger.Indent()
			//
			accessor.Min = min
			accessor.Max = max

			inner.Printf("Min : %v", min)
			inner.Printf("Max : %v", max)
		}
		return nil
	}),

	// Accessor Task
	TightPacking: FnPostTask("TightPacking", func(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
		// TODO not yet tested
		// Condition check
		// Do not access the same part of buffer from another accessor.
		// Do not access the same bufferView from another accessor.
		// Each Buffer must be cached
		for i := 0; i < len(gltf.Accessors); i++ {
			a := gltf.Accessors[i]
			_, err := a.BufferView.Buffer.Modify()
			if err != nil {
				return err
			}
			for j := i + 1; j < len(gltf.Accessors); j++ {
				b := gltf.Accessors[j]
				if a.BufferView == b.BufferView {
					return errors.Errorf("Overlap")
				}
				if a.BufferView.Buffer == b.BufferView.Buffer {
					af := a.BufferView.ByteOffset
					at := a.BufferView.ByteOffset + a.BufferView.ByteLength
					bf := b.BufferView.ByteOffset
					bt := b.BufferView.ByteOffset + b.BufferView.ByteLength

					if af >= bf{
						if bt - af <= 0{
							return errors.Errorf("Overlap")
						}
					}else{
						if at - bf <= 0{
							return errors.Errorf("Overlap")
						}
					}
				}
			}
		}
		//
		for i, v := range gltf.Accessors {
			if v.BufferView.ByteStride != 0{
				if len(v.Name) > 0{
					logger.Println("gltf.Accessor['%s']", v.Name)
				}else {
					logger.Println("gltf.Accessor[%d]", i)
				}
				inner := logger.Indent()
				var lineDataCount int
				var lineLength = v.Count
				switch v.Type {
				case SCALAR:
					lineDataCount = 1
				case VEC2:
					lineDataCount = 2
				case VEC3:
					lineDataCount = 3
				case VEC4:
					lineDataCount = 4
				case MAT2:
					lineDataCount = 2
					lineLength *= 2
				case MAT3:
					lineDataCount = 3
					lineLength *= 3
				case MAT4:
					lineDataCount = 4
					lineLength *= 4
				}
				var (
					tightBlockSize = lineDataCount * v.ComponentType.Size()
					leftBlockSize = tightBlockSize % v.BufferView.ByteStride
					blockSize = tightBlockSize + leftBlockSize
				)
				inner.Printf("Block[%d/%d/%d]", tightBlockSize, leftBlockSize, blockSize)
				if leftBlockSize == 0{
					inner.Println("No need to packing")
					v.BufferView.ByteStride = 0
				}else {
					var wrtOffset = v.BufferView.ByteOffset + v.ByteOffset
					var rdOffset = v.BufferView.ByteOffset + v.ByteOffset
					for i := 0; i < lineLength; i++ {
						copy(v.BufferView.Buffer.cache[wrtOffset:], v.BufferView.Buffer.cache[rdOffset:rdOffset + tightBlockSize])
						wrtOffset += tightBlockSize
						rdOffset += blockSize
					}
					v.BufferView.ByteStride = 0
					v.BufferView.ByteLength = wrtOffset - v.ByteOffset
					inner.Println("Packing complete")
				}
			}
		}
		return nil
	}),
	// Node Task
	ModelAlign: func(x align.Align, y align.Align, z align.Align) Task {
		return &modelAlign{x: x, y: y, z: z}
	},
	// Node Task
	ModelScale: func(axis axis.Axis, meter meter.Meter) Task {
		return &modelScale{
			len:  meter,
			axis: axis,
		}
	},
	// BufferView Task
	AutoBufferTarget: FnPreTask("Auto Buffer Target", func(parser *parserContext, gltf *SpecGLTF, logger *glog.Glogger) {
		for _, mesh := range gltf.Meshes {
			for _, prim := range mesh.Primitives {
				if prim.Indices != nil && inRange(*prim.Indices, len(gltf.Accessors)) {
					bvi := gltf.Accessors[*prim.Indices].BufferView
					if bvi != nil && inRange(*bvi, len(gltf.BufferViews)) {
						bv := &gltf.BufferViews[*bvi]
						if bv.Target != nil && *bv.Target == ELEMENT_ARRAY_BUFFER {
							if bv.Name != nil {
								logger.Printf("gltf.BufferView['%s'] Already setup : EBO", *bv.Name)
							} else {
								logger.Printf("gltf.BufferView[%d] Already setup : EBO", *bvi)
							}
						} else {
							if bv.Target == nil {
								bv.Target = new(BufferType)
								*bv.Target = ELEMENT_ARRAY_BUFFER
								if bv.Name != nil {
									logger.Printf("gltf.BufferView['%s'] Target : EBO", *bv.Name)
								} else {
									logger.Printf("gltf.BufferView[%d] Target : EBO", *bvi)
								}
							} else {
								logger.Printf("gltf.BufferView[%d] Expected EBO, but VBO detected", *bvi)
							}
						}
					}
				}
			}
		}
		for i := range gltf.BufferViews {
			if gltf.BufferViews[i].Target == nil {
				gltf.BufferViews[i].Target = new(BufferType)
				*gltf.BufferViews[i].Target = ARRAY_BUFFER
				if gltf.BufferViews[i].Name != nil {
					logger.Printf("gltf.BufferView['%s'] Target : VBO", *gltf.BufferViews[i].Name)
				} else {
					logger.Printf("gltf.BufferView[%d] Target : VBO", i)
				}
			}
		}
	}),

}

type modelAlign struct {
	x align.Align
	y align.Align
	z align.Align
}

func (s *modelAlign) TaskName() string {
	return "Model Align"
}
func (s *modelAlign) PostLoad(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
	for i, node := range gltf.Nodes {
		// TODO if node has skin
		min, max, ok := uMinMax(node)
		if !ok {
			continue
		}
		if node.Parent != nil {
			continue
		}
		if len(node.Name) > 0 {
			logger.Printf("glTF.Nodes['%s'] ", node.Name)
		} else {
			logger.Printf("glTF.Nodes[%d] ", i)
		}
		inner := logger.Indent()

		inner.Printf("Min : %v", min)
		inner.Printf("Max : %v", max)
		//
		var translate = mgl32.Vec3{
			diff(s.x, min[0], max[0]),
			diff(s.y, min[1], max[1]),
			diff(s.z, min[2], max[2]),
		}
		inner.Printf("Translate X : %v", translate[0])
		inner.Printf("Translate Y : %v", translate[1])
		inner.Printf("Translate Z : %v", translate[2])
		//uTransform(node, mgl32.Translate3D(translate[0], translate[1], translate[2]))
		if node.Matrix != mgl32.Ident4() {
			node.Matrix = node.Matrix.Mul4(mgl32.Translate3D(translate[0], translate[1], translate[2]))
		} else {
			node.Translation = node.Translation.Add(translate)
		}
		inner.Printf("Translate Complete")
	}
	return nil
}
func diff(a align.Align, min, max float32) float32 {

	const e = 0.0001
	switch a {
	case align.No:
	case align.Zero:
		if !mgl32.FloatEqualThreshold(-min, max, e) {
			return -(max + min) / 2
		}
	case align.Negative:
		if !mgl32.FloatEqualThreshold(min, 0, e) {
			return -min
		}
	case align.Positive:
		if !mgl32.FloatEqualThreshold(max, 0, e) {
			return -max
		}
	}
	return 0
}

type modelScale struct {
	axis axis.Axis
	len  meter.Meter
}

func (s *modelScale) TaskName() string {
	return "Model Scale"
}
func (s *modelScale) PostLoad(parser *parserContext, gltf *GLTF, logger *glog.Glogger) error {
	for i, node := range gltf.Nodes {
		min, max, ok := uMinMax(node)
		if !ok {
			continue
		}
		if node.Parent != nil {
			continue
		}
		if len(node.Name) > 0 {
			logger.Printf("glTF.Nodes['%s'] ", node.Name)
		} else {
			logger.Printf("glTF.Nodes[%d] ", i)
		}
		//
		inner := logger.Indent()
		inner.Printf("Min : %v", min)
		inner.Printf("Max : %v", max)
		//
		var scale float32
		switch s.axis {
		case axis.X:
			scale = s.len.Convert(prefix.No).F32() / mgl32.Abs(max.X()-min.X())
		case axis.Y:
			scale = s.len.Convert(prefix.No).F32() / mgl32.Abs(max.Y()-min.Y())
		case axis.Z:
			scale = s.len.Convert(prefix.No).F32() / mgl32.Abs(max.Z()-min.Z())

		}
		if mgl32.FloatEqualThreshold(scale, 1, 0.001) {
			inner.Printf("Scaled node")
			continue
		}
		inner.Printf("Scale : %v", scale)
		//uTransform(node, mgl32.Scale3D(scale, scale, scale))
		if node.Matrix != mgl32.Ident4() {
			node.Matrix = node.Matrix.Mul4(mgl32.Scale3D(scale, scale, scale))
		} else {
			node.Scale = mgl32.Vec3{
				node.Scale[0] * scale,
				node.Scale[1] * scale,
				node.Scale[2] * scale,
			}
		}
		inner.Printf("Scale Complete")
	}
	return nil
}


func uMinMax(node *Node) (min, max mgl32.Vec3, ok bool) {
	min = mgl32.Vec3{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max = mgl32.Vec3{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	recurMinMax(node, mgl32.Ident4(), &min, &max)
	if min[0] > max[0] {
		return min, max, false
	}
	return min, max, true
}
func recurMinMax(node *Node, mtx mgl32.Mat4, min, max *mgl32.Vec3) {
	mtx = mtx.Mul4(node.Transform())
	//
	if node.Mesh != nil {
		for _, prim := range node.Mesh.Primitives {
			posattr := prim.Attributes[POSITION]
			var (
				tempmin mgl32.Vec3
				tempmax mgl32.Vec3
			)
			if len(posattr.Min) < 0 || len(posattr.Max) < 0 {
				//TODO
				panic("TODO")
			} else {
				tempmin = mtx.Mul4x1(mgl32.Vec3{posattr.Min[0], posattr.Min[1], posattr.Min[2]}.Vec4(1)).Vec3()
				tempmax = mtx.Mul4x1(mgl32.Vec3{posattr.Max[0], posattr.Max[1], posattr.Max[2]}.Vec4(1)).Vec3()
			}
			if tempmin[0] < min[0] {
				min[0] = tempmin[0]
			}
			if tempmin[1] < min[1] {
				min[1] = tempmin[1]
			}
			if tempmin[2] < min[2] {
				min[2] = tempmin[2]
			}
			if tempmax[0] > max[0] {
				max[0] = tempmax[0]
			}
			if tempmax[1] > max[1] {
				max[1] = tempmax[1]
			}
			if tempmax[2] > max[2] {
				max[2] = tempmax[2]
			}
		}
	}
	//
	for _, child := range node.Children {
		recurMinMax(child, mtx, min, max)
	}
}
