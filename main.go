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

type Mob struct {
	AABB
	color.Color
}

func (m *Mob) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, m.X, m.Y, m.W, m.H, m.Color)
}

type Boss struct {
	Mob
	DX, DY float64
}

func (b *Boss) Update(g *Game) {
	b.DY = UpdateDelta(b.DY, ebiten.IsKeyPressed(ebiten.KeyS), ebiten.IsKeyPressed(ebiten.KeyW))
	b.DX = UpdateDelta(b.DX, ebiten.IsKeyPressed(ebiten.KeyD), ebiten.IsKeyPressed(ebiten.KeyA))

	b.X += b.DX
	b.Y += b.DY
	b.ClampToBound(ScreenW, ScreenH)

	// Boss pushes the hero, so hero does the detangle, not boss.
	g.Hero.Detangle(b.AABB)
	g.Hero.ClampToBound(ScreenW, ScreenH)
	// Hero might push back if they don't fit, so now boss detangles
	b.Detangle(g.Hero.AABB)
}

type Game struct {
	Boss *Boss
	Hero *Mob
}

func NewGame() *Game {
	boss := &Boss{
		Mob: Mob{
			AABB:  AABB{10, 10, 128, 128},
			Color: color.RGBA{0x00, 0x00, 0xFF, 0xFF},
		},
	}
	hero := &Mob{
		AABB:  AABB{250, 400, 64, 64},
		Color: color.RGBA{0x00, 0xFF, 0x00, 0xFF},
	}
	return &Game{
		boss,
		hero,
	}
}

func (g *Game) Update() error {
	g.Boss.Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xF0, 0xF0, 0xF0, 0xFF})
	g.Boss.Draw(screen)
	g.Hero.Draw(screen)
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
