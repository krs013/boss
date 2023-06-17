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

type Room struct {
	Width, Height float64
	Obstacles     []AABB
	Triggers      map[AABB]func()

	Floor  color.Color
	Wall   color.Color
	Button color.Color
}

func (r Room) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, 0, 0, r.Width, r.Height, r.Floor)
	for t := range r.Triggers {
		ebitenutil.DrawRect(dst, t.X, t.Y, t.W, t.H, r.Button)
	}
	for _, o := range r.Obstacles {
		ebitenutil.DrawRect(dst, o.X, o.Y, o.W, o.H, r.Wall)
	}
}

type Mob struct {
	AABB
	color.Color
}

func (m Mob) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, m.X, m.Y, m.W, m.H, m.Color)
}

type Boss struct {
	Mob
	DX, DY float64
}

func (b *Boss) Update(g *Game) {
	b.DX = UpdateDelta(b.DX, ebiten.IsKeyPressed(ebiten.KeyD), ebiten.IsKeyPressed(ebiten.KeyA))
	b.DY = UpdateDelta(b.DY, ebiten.IsKeyPressed(ebiten.KeyS), ebiten.IsKeyPressed(ebiten.KeyW))
	b.Move(b.DX, b.DY)

	b.DetangleRoom(g.Room)

	// Boss pushes the hero, so hero does the detangle, not boss.
	g.Hero.Detangle(b.AABB)
	g.Hero.DetangleRoom(g.Room)
	// Hero might push back if they don't fit, so now boss detangles
	b.Detangle(g.Hero.AABB)
}

type Hero struct {
	Mob
	DX, DY float64
}

func (h *Hero) Update(g *Game) {
	tx := g.Boss.X + g.Boss.W/2 - h.W/2
	ty := g.Boss.Y + g.Boss.H/2 - h.H/2

	h.DX = UpdateDelta(h.DX, tx > h.X, tx < h.X)
	h.DY = UpdateDelta(h.DY, ty > h.Y, ty < h.Y)
	h.Move(h.DX, h.DY)
}

type Game struct {
	Room *Room
	Boss *Boss
	Hero *Hero
}

func NewGame() *Game {
	room := &Room{
		Width:  ScreenW,
		Height: ScreenH,
		Obstacles: []AABB{
			{10, 10, 300, 40},
			{10, 40, 40, 200},
		},
		Triggers: map[AABB]func(){
			{500, 500, 64, 64}: nil,
		},
		Floor:  Color7,
		Wall:   Color6,
		Button: Color3,
	}
	boss := &Boss{
		Mob: Mob{
			AABB:  AABB{50, 50, 128, 128},
			Color: Color0,
		},
	}
	hero := &Hero{
		Mob: Mob{
			AABB:  AABB{250, 400, 64, 64},
			Color: Color1,
		},
	}
	return &Game{room, boss, hero}
}

func (g *Game) Update() error {
	// Boss goes after hero, since boss resolves boss-hero collision
	g.Hero.Update(g)
	g.Boss.Update(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Room.Draw(screen)
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
