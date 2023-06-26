package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// MainMenu displays the game name between hero runs.
type MainMenu struct {
	SplashFont font.Face
}

// NewMainMenu returns a new MainMenu.
func NewMainMenu() *MainMenu {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &MainMenu{
		SplashFont: font,
	}
}

// Update advances to the waiting room if any key is pressed.
func (m *MainMenu) Update(g *Game) {
	if len(ebiten.AppendInputChars(nil)) > 0 {
		g.Scene = NewWaitScene(g)
	}
}

// Draw just draws the opening splash screen.
func (m *MainMenu) Draw(screen *ebiten.Image) {
	screen.Fill(Color7)
	text.Draw(screen, "You're the Boss!!!", m.SplashFont, 400, 250, Color4)
}
