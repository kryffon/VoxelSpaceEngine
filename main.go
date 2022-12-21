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
	renderDist   = 1024

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
		plMove: true,
		pixels: make([]byte, screenHeight*screenWidth*4),
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
	if g.player.moved {
		g.plMove = true
		g.player.moved = false
	}
	return nil
}

func (g *Game) getDepth(x, y int) float64 {
	d := depthMaps[g.mapID].At(x, y)
	R, _, _, _ := d.RGBA()
	return float64(0xff * R / 0xffff)
}

func (g *Game) DrawMap() {
	height := g.player.height
	horizon := g.player.horizon
	scale_height := 240.0

	ybuffer := [screenWidth]float64{}
	for i := range ybuffer {
		ybuffer[i] = screenHeight
	}

	psin, pcos := math.Sincos(g.player.phi)
	dz := 1.0
	for z := 1.0; z <= renderDist; z += dz {
		plx := -pcos*z - psin*z
		ply := psin*z - pcos*z
		prx := pcos*z - psin*z
		pry := -psin*z - pcos*z

		dx := (prx - plx) / screenWidth
		dy := (pry - ply) / screenHeight
		// add player coords after as it doesn't matter. they cancel out
		plx += g.player.x
		ply += g.player.y
		for i := 0; i < screenWidth; i++ {
			// for repeating map
			ix, iy := int(plx)&(mapHeight-1), int(ply)&(mapWidth-1)

			depth := g.getDepth(ix, iy)
			heightOnScreen := horizon + (height-depth)/z*scale_height

			c := colorMaps[g.mapID].At(ix, iy)
			g.DrawVerticalLine(i, int(heightOnScreen), int(ybuffer[i]), c)
			if heightOnScreen <= ybuffer[i] {
				ybuffer[i] = heightOnScreen
			}
			plx += dx
			ply += dy
		}
		dz += 0.005
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
			g.pixels[i] = 0xe0
		}
		log.Println("DRAW - START - MAP:", g.mapID)
		g.DrawMap()
		log.Println("DRAW - COMPLETE - MAP:", g.mapID)
		g.plMove = false
	}
	// screen.Fill(color.RGBA{245, 245, 245, 245})
	screen.WritePixels(g.pixels)

	msg := fmt.Sprintf("FPS: %0.0f, TPS: %0.0f\nMap(%d): MN\nMOVE: WASD\n Pitch: QE\nHeight: RF",
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
