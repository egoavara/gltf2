package gltf2

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func inRange(testset SpecGLTFID, len int) bool {
	return 0 <= int(testset) && int(testset) < len
}

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
