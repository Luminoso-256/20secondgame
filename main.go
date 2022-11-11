package main

import (
	"embed"
	_ "image/png"
	"log"

	"minigame/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed data/*
var embedFS embed.FS

const (
	Dbg_LOCAL_ASSETS   = false
	Dbg_NO_TITLESCREEN = false
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SaltRush")
	s := 0
	if Dbg_NO_TITLESCREEN {
		s = 1
	}
	g := &game.Game{
		Player: game.Player{
			X: 640 / 2, Y: 480 / 2,
		},
		Assets:    game.LoadAssets(Dbg_LOCAL_ASSETS, embedFS),
		GameState: s,
	}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
