package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawLightCircle(dst *ebiten.Image, cx, cy, r float64, clr color.Color) {
	var path vector.Path
	rd, g, b, a := clr.RGBA()

	path.Arc(float32(cx), float32(cy), float32(r), 0, 2*math.Pi, vector.Clockwise)

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = float32(rd) / 0xffff
		vertices[i].ColorG = float32(g) / 0xffff
		vertices[i].ColorB = float32(b) / 0xffff
		vertices[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.CompositeMode = ebiten.CompositeModeLighter

	dst.DrawTriangles(vertices, indices, emptySubImage, nil)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func roundUp(numToRound, multiple int) int {
	if multiple == 0 {
		return numToRound
	}

	remainder := Abs(numToRound) % multiple
	if remainder == 0 {
		return numToRound
	}

	if numToRound < 0 {
		return -(Abs(numToRound) - remainder)
	} else {
		return numToRound + multiple - remainder
	}
}
