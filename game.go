package submarinegame

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gitlab.com/high-creek-software/go2d/loader"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log/slog"
	"os"
	"time"
)

var SCREEN_WIDTH = 1280.0
var SCREEN_HEIGHT = 720.0
var WATER_SURFACE = 140.0

type SoundType int

const (
	SOUNDS_SPLASH SoundType = iota
	SOUNDS_TORPEDO
	SOUNDS_HIT
	SOUNDS_EXP_UNDERWATER
	SOUNDS_MISSLE
	SOUNDS_PING
)

var mplusFaceSource *text.GoTextFaceSource
var gamePrinter *message.Printer

//go:embed assets
var assets embed.FS
var assetLoader *loader.AssetLoader

type Game struct {
	ship               *Ship
	levelIndex         int
	level              *Level
	startScreen        *StartScreen
	startUI            *StartUI
	gameOverScreen     *GameOverScreen
	pointerX, pointerY int

	audioContext  *audio.Context
	seaLoop       *audio.Player
	splashSound   *audio.Player
	expUnderwater *audio.Player
	missle        *audio.Player
	hit           *audio.Player
	ping          *audio.Player
}

func NewGame() *Game {
	subGame := &Game{}
	subGame.initializeLoader()
	subGame.initializeUIComponents()
	subGame.initializeSounds()

	subGame.startScreen = NewStartScreen(subGame.GameStarted)
	//subGame.startUI = NewStartUI()

	return subGame
}

func (g *Game) Update() error {
	if g.level != nil {
		g.level.Update()
	} else if g.gameOverScreen != nil {
		g.gameOverScreen.Update()
	} else if g.startScreen != nil {
		g.startScreen.Update()
	} else if g.startUI != nil {
		g.startUI.Update()
	}

	g.pointerX, g.pointerY = ebiten.CursorPosition()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.level != nil {
		g.level.Draw(screen)
	} else if g.gameOverScreen != nil {
		g.gameOverScreen.Draw(screen)
	} else if g.startScreen != nil {
		g.startScreen.Draw(screen)
	} else if g.startUI != nil {
		g.startUI.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d, %d", g.pointerX, g.pointerY))
}

func (g *Game) Layout(ow, oh int) (w, h int) {
	SCREEN_WIDTH = float64(ow)
	SCREEN_HEIGHT = float64(oh)
	return ow, oh
}

func (g *Game) GameStarted(shipType ShipType) {
	g.ship = NewShip(shipType)
	g.newLevel()
}

func (g *Game) Restart() {
	g.gameOverScreen = nil
	g.startScreen = NewStartScreen(g.GameStarted)
}

func (g *Game) LevelComplete() {
	g.newLevel()
}

func (g *Game) newLevel() {
	g.level = NewLevel(g, g.levelIndex, g.ship)
	g.levelIndex += 1
}

func (g *Game) GameOver() {
	g.level = nil
	g.gameOverScreen = NewGameOverScreen(g.Restart)
	g.levelIndex = 0
}

func (g *Game) PlaySound(soundType SoundType) {
	switch soundType {
	case SOUNDS_SPLASH:
		if g.splashSound != nil {
			g.splashSound.Rewind()
			g.splashSound.Play()
		}
	case SOUNDS_EXP_UNDERWATER:
		if g.expUnderwater != nil {
			g.expUnderwater.Rewind()
			g.expUnderwater.Play()
		}
	case SOUNDS_MISSLE:
		if g.missle != nil {
			g.missle.Rewind()
			g.missle.Play()
		}
	case SOUNDS_HIT:
		if g.hit != nil {
			g.hit.Rewind()
			g.hit.Play()
		}
	case SOUNDS_PING:
		if g.ping != nil {
			g.ping.Rewind()
			g.ping.Play()
		}
	}
}

func (g *Game) PlaySoundDelay(soundType SoundType, delay time.Duration) {
	go func() {
		<-time.Tick(delay)
		g.PlaySound(soundType)
	}()
}

func (g *Game) initializeUIComponents() {
	var err error
	mplusFaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		slog.Error("error loading font", "error", err)
		os.Exit(2)
	}

	gamePrinter = message.NewPrinter(language.English)
}

func (g *Game) initializeLoader() {
	assetLoader = loader.NewAssetLoader(assets)
}

func (g *Game) initializeSounds() {
	g.audioContext = audio.NewContext(48_000)

	var err error
	g.splashSound, err = g.createAudioPlayer("assets/sounds/splash.mp3")
	if err != nil {
		slog.Info("error creating splash", "error", err)
	} else {
		g.splashSound.SetVolume(0.33)
	}
	g.expUnderwater, err = g.createAudioPlayer("assets/sounds/exp_underwater.mp3")
	if err != nil {
		slog.Info("error creating underwater subexplosion", "error", err)
	} else {
		g.expUnderwater.SetVolume(0.22)
	}
	g.missle, err = g.createAudioPlayer("assets/sounds/missle.mp3")
	if err != nil {
		slog.Info("error creating missle", "error", err)
	} else {
		g.missle.SetVolume(0.33)
	}
	g.hit, err = g.createAudioPlayer("assets/sounds/hit.mp3")
	if err != nil {
		slog.Info("error creating hit", "error", err)
	} else {
		g.hit.SetVolume(0.22)
	}
	g.ping, err = g.createAudioPlayer("assets/sounds/ping.mp3")
	if err != nil {
		slog.Info("error creating ping", "error", err)
	} else {
		g.ping.SetVolume(0.22)
	}

	if seaLoopReader, err := assetLoader.GetReader("assets/sounds/sea.mp3"); err == nil {
		if stream, streamErr := mp3.DecodeWithoutResampling(seaLoopReader); streamErr == nil {
			loop := audio.NewInfiniteLoop(stream, 24952000)
			g.seaLoop, _ = g.audioContext.NewPlayer(loop)
			g.seaLoop.SetVolume(0.15)
			g.seaLoop.Play()
		}
	}

}

func (g *Game) createAudioPlayer(path string) (*audio.Player, error) {
	if r, err := assetLoader.GetReader(path); err == nil {
		if stream, streamErr := mp3.DecodeWithoutResampling(r); streamErr == nil {
			return g.audioContext.NewPlayer(stream)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
