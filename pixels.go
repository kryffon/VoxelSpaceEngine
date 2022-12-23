package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	for z := 1.0; z <= float64(g.renderDist); z += dz {
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

func (g *Game) DrawPixels(screen *ebiten.Image) {
	if g.plMove {
		for i := range g.pixels {
			g.pixels[i] = 0xe0
		}
		log.Println("DRAW - PIX - START - MAP:", g.mapID)
		g.DrawMap()
		log.Println("DRAW - PIX - COMPLETE - MAP:", g.mapID)
		g.plMove = false
	}
	screen.WritePixels(g.pixels)
}
