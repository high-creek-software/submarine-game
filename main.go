package main

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Person struct {
	Name   string
	Age    int
	Height float64
}

func (person Person) String() string {
	return person.Name
}

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(ow, oh int) (w, h int) {

	return 1280, 720
}

func main() {

	// p := Person{
	// 	Name:   "John Doe",
	// 	Age:    76,
	// 	Height: 36.5,
	// }

	// g := Person{
	// 	Name:   "George",
	// 	Age:    32,
	// 	Height: 45,
	// }

	// slog.Info("George", "name", g.String())

	subGame := &Game{}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Submarine Game")

	ebiten.RunGame(subGame)
}

func printPerson(p1, p2 Person) {
	slog.Info("Person", "p1", p1, "p2", p2)
}
