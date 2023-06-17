package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenW = 1920
	ScreenH = 1080
)

var hero = AABB{250, 400, 64, 64}

type Boss struct {
	AABB
	DX, DY float64
	color.Color
}

func (b *Boss) Update() {
	b.DY = UpdateDelta(b.DY, ebiten.IsKeyPressed(ebiten.KeyS), ebiten.IsKeyPressed(ebiten.KeyW))
	b.DX = UpdateDelta(b.DX, ebiten.IsKeyPressed(ebiten.KeyD), ebiten.IsKeyPressed(ebiten.KeyA))

	b.X += b.DX
	b.Y += b.DY

	b.Detangle(hero)
}

func (b *Boss) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, b.X, b.Y, b.W, b.H, b.Color)
}

type Game struct {
	Boss *Boss
}

func NewGame() *Game {
	boss := &Boss{
		AABB:  AABB{10, 10, 128, 128},
		Color: color.RGBA{0x00, 0x00, 0xFF, 0xFF},
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

	c := color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	if g.Boss.Collide(hero) {
		c = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	}
	ebitenutil.DrawRect(screen, hero.X, hero.Y, hero.W, hero.H, c)
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
