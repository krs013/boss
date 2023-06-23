package main

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

// NYX8 Palete - https://lospec.com/palette-list/nyx8
var (
	Color0 = color.RGBA{0xF6, 0xD6, 0xBD, 0xFF}
	Color1 = color.RGBA{0xC3, 0xA3, 0x8A, 0xFF}
	Color2 = color.RGBA{0x99, 0x75, 0x77, 0xFF}
	Color3 = color.RGBA{0x81, 0x62, 0x71, 0xFF}
	Color4 = color.RGBA{0x4E, 0x49, 0x5F, 0xFF}
	Color5 = color.RGBA{0x20, 0x39, 0x4F, 0xFF}
	Color6 = color.RGBA{0x0F, 0x2A, 0x3F, 0xFF}
	Color7 = color.RGBA{0x08, 0x14, 0x1E, 0xFF}
)

type Animation struct {
	Sheet            *ebiten.Image
	OffsetX, OffsetY int
	Width, Height    int
	NumFrames        int
	FrameSpeed       int
	Transform        ebiten.GeoM

	FrameCount int
}

func (a *Animation) Sprite() *ebiten.Image {
	a.FrameCount += 1
	i := (a.FrameCount / a.FrameSpeed) % a.NumFrames
	sx := a.OffsetX + i*a.Width
	rect := image.Rect(sx, a.OffsetY, sx+a.Width, a.OffsetY+a.Height)
	return a.Sheet.SubImage(rect).(*ebiten.Image)
}

func (a *Animation) Draw(dst *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Concat(a.Transform)
	op.GeoM.Translate(x, y)
	dst.DrawImage(a.Sprite(), op)
}

func BossAnimations() (idle, left, right *Animation) {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	sheet := ebiten.NewImageFromImage(img)

	op := ebiten.GeoM{}
	op.Scale(4, 4)
	idle = &Animation{
		Sheet:      sheet,
		OffsetX:    0,
		OffsetY:    0,
		Width:      32,
		Height:     32,
		NumFrames:  5,
		FrameSpeed: 5,
		Transform:  op,
	}

	right = &Animation{
		Sheet:      sheet,
		OffsetX:    0,
		OffsetY:    32,
		Width:      32,
		Height:     32,
		NumFrames:  5,
		FrameSpeed: 5,
		Transform:  op,
	}

	op.Scale(-1, 1)
	op.Translate(128, 0)
	left = &Animation{
		Sheet:      sheet,
		OffsetX:    0,
		OffsetY:    32,
		Width:      32,
		Height:     32,
		NumFrames:  5,
		FrameSpeed: 5,
		Transform:  op,
	}

	return idle, left, right
}
