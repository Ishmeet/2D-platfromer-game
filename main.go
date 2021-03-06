package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	rplatformer "github.com/hajimehoshi/ebiten/examples/resources/images/platformer"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	// Settings
	screenWidth  = 640
	screenHeight = 448
)

var (
	leftSprite      *ebiten.Image
	rightSprite     *ebiten.Image
	idleSprite      *ebiten.Image
	backgroundImage *ebiten.Image
)

func init() {
	// Preload images
	img, _, err := image.Decode(bytes.NewReader(rplatformer.Right_png))
	if err != nil {
		panic(err)
	}
	rightSprite, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(rplatformer.Left_png))
	if err != nil {
		panic(err)
	}
	leftSprite, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, _, err = image.Decode(bytes.NewReader(rplatformer.MainChar_png))
	if err != nil {
		panic(err)
	}
	idleSprite, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	backgroundImage, _, err = ebitenutil.NewImageFromFile("Mario.png", ebiten.FilterDefault)

	// img, _, err = image.Decode(bytes.NewReader(rplatformer.Background_png))
	if err != nil {
		panic(err)
	}
	// backgroundImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

const (
	unit    = 16
	groundY = 380
)

type char struct {
	x  int
	y  int
	vx int
	vy int
}

func (c *char) tryJump() {
	// Now the character can jump anytime, even when the character is not on the ground.
	// If you want to restrict the character to jump only when it is on the ground, you can add an 'if' clause:
	//
	//     if gopher.y == groundY * unit {
	//         ...
	c.vy = -10 * unit
}

func (c *char) update() {
	c.x += c.vx
	c.y += c.vy
	if c.y > groundY*unit {
		c.y = groundY * unit
	}
	if c.vx > 0 {
		c.vx -= 4
	} else if c.vx < 0 {
		c.vx += 4
	}
	if c.vy < 20*unit {
		c.vy += 8
	}
}

func (c *char) draw(screen *ebiten.Image) {
	s := idleSprite
	switch {
	case c.vx > 0:
		s = rightSprite
	case c.vx < 0:
		s = leftSprite
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(float64(c.x)/unit, float64(c.y)/unit)
	screen.DrawImage(s, op)
}

// Game ...
type Game struct {
	gopher *char
}

// Update ...
func (g *Game) Update(screen *ebiten.Image) error {
	if g.gopher == nil {
		g.gopher = &char{x: 200 * unit, y: groundY * unit}
	}

	// Controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// g.gopher.vx = -4 * unit
		pos += 4
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		// g.gopher.vx = 4 * unit
		pos -= 4
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.gopher.tryJump()
	}
	g.gopher.update()
	return nil
}

var pos float64

// Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	// Draws Background Image
	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(pos, 0)
	screen.DrawImage(backgroundImage, op)

	// Draws the Gopher
	ebitenutil.DrawRect(screen, float64(g.gopher.x)/unit, float64(g.gopher.y)/unit, 20, 20, color.RGBA{0xff, 0x00, 0x00, 0xff})
	// g.gopher.draw(screen)

	// Show the message
	msg := fmt.Sprintf("TPS: %0.2f\nPress the space key to jump.", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Platformer (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
