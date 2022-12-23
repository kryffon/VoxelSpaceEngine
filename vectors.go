package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(5, 5)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

type Line struct {
	xs  float32
	ys  float32
	xd  float32
	yd  float32
	col color.Color
}

func (g *Game) GetLinesFromMap() {
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
			g.AddLine(i, int(heightOnScreen), int(ybuffer[i]), c)
			if heightOnScreen <= ybuffer[i] {
				ybuffer[i] = heightOnScreen
			}
			plx += dx
			ply += dy
		}
		dz += 0.005
	}
}

func (g *Game) AddLine(x, ytop, ybot int, col color.Color) {
	if ytop < 0 {
		ytop = 0
	}
	if ytop > ybot {
		return
	}
	g.lines = append(g.lines, Line{
		xs:  float32(x),
		ys:  float32(ytop),
		xd:  float32(x),
		yd:  float32(ybot),
		col: col,
	})
}

func (g *Game) DrawAllLines(screen *ebiten.Image) {
	for _, vec := range g.lines {
		var path vector.Path
		path.MoveTo(vec.xs, vec.ys)
		path.LineTo(vec.xd, vec.yd)

		op := &vector.StrokeOptions{}
		op.LineCap = vector.LineCapSquare
		op.LineJoin = vector.LineJoinMiter
		op.Width = 2
		vs, is := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, op)
		R, G, B, A := vec.col.RGBA()
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = float32(R) / 65535.0
			vs[i].ColorG = float32(G) / 65535.0
			vs[i].ColorB = float32(B) / 65535.0
			vs[i].ColorA = float32(A) / 65535.0
		}
		screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{})
	}
}

func (g *Game) DrawVectors(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 245, 245, 245})
	if g.plMove {
		g.lines = []Line{}
		log.Println("CALC - VEC - START - MAP:", g.mapID)
		g.GetLinesFromMap()
		log.Println("CALC - VEC - COMPLETE - MAP:", g.mapID)
		g.plMove = false
	}
	g.DrawAllLines(screen)
}
