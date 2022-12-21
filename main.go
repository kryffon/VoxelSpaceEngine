package main

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 800
	imgWidth     = 1024
	imgHeight    = 1024

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
	mapID  int
	player Player
	plMove bool
	pixels []byte
}

type Player struct {
	x float64
	y float64
	// phi float64
}

func NewGame() *Game {
	return &Game{
		mapID:  0,
		player: Player{x: 512, y: 512},
		plMove: false,
		pixels: make([]byte, screenHeight*screenWidth*4),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.mapID = (mapCount + g.mapID - 1) % mapCount
		g.plMove = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.mapID = (g.mapID + 1) % mapCount
		g.plMove = true
	}
	return nil
}

func (g *Game) getDepth(x, y float64) float64 {
	d := depthMaps[g.mapID].At(int(x), int(y))
	R, _, _, _ := d.RGBA()
	return float64(0xff * R / 0xffff)
}

func (g *Game) DrawMap() {
	height := float64(255 * 5 / 5)
	horizon := float64(screenHeight / 3)
	scale_height := 240.0
	distance := float64(imgHeight * 2.2 / 5)

	ybuffer := [screenWidth]float64{}
	for i := range ybuffer {
		ybuffer[i] = screenHeight
	}

	dz := 0.01
	for z := 1.0; z <= distance; z += dz {
		plx, ply := -z+g.player.x, z+g.player.y
		prx := z + g.player.x
		dx := (prx - plx) / screenHeight
		// log.Println(z, int(plx), ply, prx, dx)
		for i := 0; i < screenWidth; i++ {
			depth := g.getDepth(plx, ply)
			heightOnScreen := horizon + (height-depth)/z*scale_height

			c := colorMaps[g.mapID].At(int(plx), int(ply))
			g.DrawVerticalLine(i, int(heightOnScreen), int(ybuffer[i]), c)
			if heightOnScreen <= ybuffer[i] {
				ybuffer[i] = heightOnScreen
			}
			plx += dx
		}
		dz += 0.01
	}
}

func (g *Game) DrawVerticalLine(x, ytop, ybot int, col color.Color) {
	if ytop < 0 {
		ytop = 0
	}
	if ytop > ybot {
		return
	}
	R, G, B, A := col.RGBA()
	for y := ytop; y < ybot; y++ {
		i := y*screenWidth + x
		g.pixels[4*i] = uint32ToByte(R)
		g.pixels[4*i+1] = uint32ToByte(G)
		g.pixels[4*i+2] = uint32ToByte(B)
		g.pixels[4*i+3] = uint32ToByte(A)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.plMove {

		for i := range g.pixels {
			g.pixels[i] = 0xf0
		}
		log.Println("DRAW - START - MAP:", g.mapID)
		g.DrawMap()
		log.Println("DRAW - COMPLETE - MAP:", g.mapID)
		g.plMove = false
	}
	// screen.Fill(color.RGBA{245, 245, 245, 245})
	screen.WritePixels(g.pixels)

	msg := fmt.Sprintf("FPS: %0.0f, TPS: %0.0f\nMap(%d): QE",
		ebiten.ActualFPS(), ebiten.ActualTPS(), g.mapID)
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
