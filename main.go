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

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("20-Second Game")
	g := &game.Game{
		Player: game.Player{
			X: 640 / 2, Y: 480 / 2,
		},
		Assets: game.LoadAssets(true, embedFS),
	}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
