package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Menu struct {
	SplashFont font.Face
}

func NewMenu() *Menu {
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

	return &Menu{
		SplashFont: font,
	}
}

func (m *Menu) Update(g *Game) error {
	if len(ebiten.AppendInputChars(nil)) > 0 {
		g.Scene = NewWaitScene()
	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(Color7)
	text.Draw(screen, "You're the Boss!!", m.SplashFont, 400, 250, Color4)
}
