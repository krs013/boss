package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
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

type Trigger struct {
	AABB
	Fn func()
}

type Room struct {
	Width, Height float64
	Obstacles     []AABB
	Triggers      []Trigger

	Floor  color.Color
	Wall   color.Color
	Button color.Color
}

func (r Room) Update(g *Game) {
	for _, t := range r.Triggers {
		if g.Boss.Collide(t.AABB) {
			if t.Fn != nil {
				t.Fn()
			}
		}
	}
}

func (r Room) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, 0, 0, r.Width, r.Height, r.Floor)
	for _, t := range r.Triggers {
		ebitenutil.DrawRect(dst, t.X, t.Y, t.W, t.H, r.Button)
	}
	for _, o := range r.Obstacles {
		ebitenutil.DrawRect(dst, o.X, o.Y, o.W, o.H, r.Wall)
	}
}

type Mob struct {
	AABB
	TargetX, TargetY float64
	DX, DY           float64

	color.Color
}

func (m Mob) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, m.X, m.Y, m.W, m.H, m.Color)
}

// TrackTarget updates the Mob DX, DY in the direction of the target.
func (m *Mob) TrackTarget() {
	m.DX = UpdateDelta(m.DX, m.X, m.TargetX)
	m.DY = UpdateDelta(m.DY, m.Y, m.TargetY)
}

func (m *Mob) Move() {
	m.Translate(m.DX, m.DY)
}

type Boss struct {
	Mob
}

func (b *Boss) KeyTarget() (x, y float64) {
	x, y = b.X, b.Y
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		x += ScreenWidth
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		x -= ScreenWidth
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		y += ScreenHeight
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		y -= ScreenHeight
	}
	return x, y
}

func (b *Boss) Update(g *Game) {
	tx, ty := b.KeyTarget()
	b.DX = UpdateDelta(b.DX, b.X, tx)
	b.DY = UpdateDelta(b.DY, b.Y, ty)

	b.Move()
	b.DetangleRoom(g.Room)
	b.ClampToRoom(g.Room)

	// Boss pushes the hero, so hero does the detangle, not boss.
	g.Hero.Detangle(b.AABB)
	// After hero gets pushed, hero might need to avoid obstacles.
	g.Hero.DetangleRoom(g.Room)
	g.Hero.ClampToRoom(g.Room)
	// Hero might push back if they don't fit, so now boss detangles
	b.Detangle(g.Hero.AABB)
}

type Hero struct {
	Mob
}

func (h *Hero) Update(g *Game) {
	// Set target so hero center moves towards the boss center.
	h.TargetX = g.Boss.X + g.Boss.W/2 - h.W/2
	h.TargetY = g.Boss.Y + g.Boss.H/2 - h.H/2
	h.TargetX = 800
	h.TargetY = 300
	h.TrackTarget()

	// Just do the move. Boss will handle collisions and pushing.
	h.Move()
}

type Game struct {
	Room *Room
	Boss *Boss
	Hero *Hero
}

func NewGame() *Game {
	room := &Room{
		Width:  ScreenWidth,
		Height: ScreenHeight,
		Obstacles: []AABB{
			{10, 10, 300, 40},
			{10, 40, 40, 200},
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
	room.Triggers = append(room.Triggers, Trigger{
		AABB: AABB{500, 500, 64, 64},
		Fn: func() {
			hero.X = ScreenWidth / 2
			hero.Y = ScreenHeight / 2
		},
	})
	return &Game{room, boss, hero}
}

func (g *Game) Update() error {
	// NB: Order matters here! Only the Boss resolves hero-boss pushing
	// interaction, so boss must go after hero has done naive moves.
	g.Hero.Update(g)
	g.Boss.Update(g)
	// Once boss has moved, see if any triggers were tripped.
	g.Room.Update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// NB: Order matters here! Later stuff draws over earlier stuff.
	g.Room.Draw(screen)
	g.Hero.Draw(screen)
	g.Boss.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Don't get fancy here - just let Ebitengine handle scaling stuff for us.
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("You're the Boss!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
