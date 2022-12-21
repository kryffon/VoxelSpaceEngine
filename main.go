package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1024
	screenHeight = 1024

	mapCount  = 14
	assetsDir = `.\maps\`
)

var (
	colorMaps []*ebiten.Image
	depthMaps []*ebiten.Image
)

func init() {
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
	mapID   int
	mapType bool
	img     *ebiten.Image
}

func (g *Game) loadImage() {
	if g.mapType {
		g.img = colorMaps[g.mapID]
	} else {
		g.img = depthMaps[g.mapID]
	}
	// log.Println("DEBUG: g ", g.mapID, g.mapType)
}

func NewGame() *Game {
	g := &Game{
		mapID:   0,
		mapType: true,
	}
	g.loadImage()
	return g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.mapID = (mapCount + g.mapID - 1) % mapCount
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.mapID = (g.mapID + 1) % mapCount
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.mapType = !g.mapType
	}
	g.loadImage()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.img, op)

	msg := fmt.Sprintf("FPS: %0.0f, TPS: %0.0f\nChange Map: QE\nChange Type: T",
		ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("VoxelSpaceEngine")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
