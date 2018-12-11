package gltf2

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(perspectiveInfinitable(mgl32.DegToRad(45), 1920.0/1080.0, 0.5, float32(math.Inf(1))))
	fmt.Println(mgl32.Perspective(mgl32.DegToRad(45), 1920.0/1080.0, 0.5, float32(math.Inf(1))))
}
