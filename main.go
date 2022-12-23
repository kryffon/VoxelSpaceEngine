package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 800
	mapWidth     = 1024
	mapHeight    = 1024

	mapCount  = 14
	assetsDir = `.\maps\`
)

var (
	colorMaps []*ebiten.Image
	depthMaps []*ebiten.Image
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	filenames, err := getFileNames(assetsDir)
	if err != nil {
		panic(err)
	}
	colorMaps = make([]*ebiten.Image, mapCount)
	depthMaps = make([]*ebiten.Image, mapCount)
	re := regexp.MustCompile("[0-9]+")
	for _, file := range filenames {
		i, _ := strconv.Atoi(re.FindString(file))
		if i <= mapCount {
			img, err := loadImage(assetsDir + file)
			if err != nil {
				panic(err)
			}
			if file[0] == 'C' {
				colorMaps[i-1] = img
			} else {
				depthMaps[i-1] = img
			}
		}
	}
}

type Game struct {
	mapID      int
	player     Player
	plMove     bool
	pixels     []byte
	lines      []Line
	renderDist int
	renderType bool
}

func NewGame() *Game {
	return &Game{
		mapID: 0,
		player: Player{
			x:       512,
			y:       512,
			phi:     math.Pi / 4,
			height:  255,
			horizon: float64(screenHeight / 2),
		},
		plMove:     true,
		pixels:     make([]byte, screenHeight*screenWidth*4),
		lines:      []Line{},
		renderDist: 500,
		renderType: true,
	}
}

func (g *Game) UpdateRenderDist(dr int) {
	if g.renderDist+dr > 0 && g.renderDist+dr <= 1000 {
		g.renderDist += dr
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.mapID = (mapCount + g.mapID - 1) % mapCount
		g.plMove = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.mapID = (g.mapID + 1) % mapCount
		g.plMove = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.renderType = !g.renderType
		g.plMove = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Move(10.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Move(-9.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Rotate(0.03)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Rotate(-0.03)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.player.ChangePitch(10.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.player.ChangePitch(-9.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.player.ChangeHeight(10.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.player.ChangeHeight(-9.0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.UpdateRenderDist(-50)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.UpdateRenderDist(50)
	}
	if g.player.moved {
		g.plMove = true
		g.player.moved = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.renderType {
		g.DrawPixels(screen)
	} else {
		g.DrawVectors(screen)
	}

	t := "pixels"
	if !g.renderType {
		t = "vectors"
	}
	msg := fmt.Sprintf("FPS: %0.0f, TPS: %0.0f\n"+
		"Map(%d): MN\nMOVE: WASD\n Pitch: QE\nHeight: RF\n"+
		"RenderDist(%d) CV\nRenderType(%s) P",
		ebiten.ActualFPS(), ebiten.ActualTPS(), g.mapID, g.renderDist, t)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	whiteSubImage.Fill(color.White)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("VoxelSpaceEngine")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
