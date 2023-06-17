package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenW = 1920
	ScreenH = 1080

	DeltaDecay = .75
	DeltaStep  = 1
	DeltaMax   = 5
)

type AABB struct {
	X, Y float64
	W, H float64
}

type Boss struct {
	AABB
	DX, DY float64
	color.Color

	Keys []ebiten.Key
}

func UpdateDelta(delta float64, inc, dec bool) float64 {
	if inc == dec {
		return delta * DeltaDecay
	}
	if inc {
		delta += DeltaStep
	}
	if dec {
		delta -= DeltaStep
	}
	return Clamp(-DeltaMax, delta, DeltaMax)
}

func (b *Boss) Update() {
	b.Keys = inpututil.AppendPressedKeys(b.Keys[:0])
	b.DY = UpdateDelta(b.DY, SliceContains(b.Keys, ebiten.KeyS), SliceContains(b.Keys, ebiten.KeyW))
	b.DX = UpdateDelta(b.DX, SliceContains(b.Keys, ebiten.KeyD), SliceContains(b.Keys, ebiten.KeyA))

	b.X = Clamp(0, b.X+b.DX, ScreenW-b.W)
	b.Y = Clamp(0, b.Y+b.DY, ScreenH-b.H)
}

func (b *Boss) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, b.X, b.Y, b.W, b.H, b.Color)
}

type Game struct {
	Boss *Boss
}

func NewGame() *Game {
	boss := &Boss{
		AABB:  AABB{10, 10, 64, 64},
		Color: color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	}
	return &Game{
		boss,
	}
}

func (g *Game) Update() error {
	g.Boss.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xF0, 0xF0, 0xF0, 0xFF})
	g.Boss.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenW, ScreenH
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
