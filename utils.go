package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func inRange(testset SpecGLTFID, len int) bool {
	return 0 <= int(testset) && int(testset) < len
}

//
//func loadBuffer(dir http.Dir, buf Buffer) ([]byte, Buffer, error) {
//	switch tbuf := buf.(type) {
//	case *LoadedBuffer:
//		return tbuf.Buffer, tbuf.Original, nil
//	case *URIBuffer:
//		f, err := dir.Open(tbuf.URI.Path)
//		if err != nil {
//			return nil, nil, err
//		}
//		defer f.Close()
//		bts, err := ioutil.ReadAll(f)
//		if err != nil {
//			return nil, nil, err
//		}
//		return bts[:tbuf.ByteLength], tbuf, nil
//	}
//	return nil, nil, errors.New("Unexpected Buffer type")
//}
//
//func loadImage(dir http.Dir, gltf *GLTF, v Image) (image.Image, Image, error) {
//	switch img := v.(type) {
//	case *LoadedImage:
//		return img.Image, img.Original, nil
//	case *BufferViewImage:
//		if !inRange(img.BufferView, len(gltf.BufferViews)) {
//			return nil, nil, errors.WithMessage(ErrorGLTFLink, "Cant find gltf.BufferViews glTFid")
//		}
//		view := gltf.BufferViews[img.BufferView]
//		if !inRange(view.Buffer, len(gltf.Buffers)) {
//			return nil, nil, errors.WithMessage(ErrorGLTFLink, "Cant find gltf.Buffers glTFid")
//		}
//
//		switch buf := gltf.Buffers[view.Buffer].(type) {
//		case *LoadedBuffer:
//			if view.ByteStride != 0 {
//				// TODO : Support ByteStride
//				return nil, nil, errors.New("Unsupport ByteStride")
//			}
//			switch img.MimeType {
//			case ImageJPEG:
//				ld, err := jpeg.Decode(bytes.NewBuffer(buf.Buffer[view.ByteOffset : view.ByteOffset+view.ByteLength]))
//				if err != nil {
//					return nil, nil, err
//				}
//				return ld, buf, nil
//			case ImagePNG:
//				ld, err := png.Decode(bytes.NewBuffer(buf.Buffer[view.ByteOffset : view.ByteOffset+view.ByteLength]))
//				if err != nil {
//					return nil, nil, err
//				}
//				return ld, buf, nil
//			default:
//				return nil, nil, errors.New("Unknown Mime '" + img.MimeType.String() + "'")
//			}
//		case *URIBuffer:
//			f, err := dir.Open(buf.URI.Path)
//			if err != nil {
//				return nil, nil, err
//			}
//			defer f.Close()
//			bts, err := ioutil.ReadAll(f)
//			if err != nil {
//				return nil, nil, err
//			}
//			ld, _, err := image.Decode(bytes.NewBuffer(bts[:buf.ByteLength]))
//			if err != nil {
//				return nil, nil, err
//			}
//			return ld, buf, nil
//		}
//	case *URIImage:
//		f, err := dir.Open(img.URI.Path)
//		if err != nil {
//			return nil, nil, err
//		}
//		defer f.Close()
//		ld, _, err := image.Decode(f)
//		if err != nil {
//			return nil, nil, err
//		}
//		return ld, img, nil
//	}
//	return nil, nil, errors.New("Unexpected Images type")
//}

func perspectiveInfinitable(fovy, aspect, near, far float32) mgl32.Mat4 {
	if math.IsInf(float64(far), 1) {
		// https://social.msdn.microsoft.com/Forums/en-US/4e7a09f4-d4b1-4251-8033-db33d5756ebd/infinite-projection-matrix?forum=xnaframework
		// fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
		const e = 0.000001
		f := float32(1. / math.Tan(float64(fovy)/2.0))
		return mgl32.Mat4{f / aspect, 0, 0, 0, 0, f, 0, 0, 0, 0, -1 + e, -1, 0, 0, -near, 0}
	}
	return mgl32.Perspective(fovy, aspect, near, far)
}

func isValidF32Color3(c mgl32.Vec3) bool {
	if c[0] < 0 || c[0] > 1.0 {
		return false
	}
	if c[1] < 0 || c[1] > 1.0 {
		return false
	}
	if c[2] < 0 || c[2] > 1.0 {
		return false
	}
	return true
}

func isValidF32Color4(c mgl32.Vec4) bool {
	if c[0] < 0 || c[0] > 1.0 {
		return false
	}
	if c[1] < 0 || c[1] > 1.0 {
		return false
	}
	if c[2] < 0 || c[2] > 1.0 {
		return false
	}
	if c[3] < 0 || c[3] > 1.0 {
		return false
	}
	return true
}
