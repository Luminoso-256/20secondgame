package game

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	USE_TIMECHECK = true
)

var (
	DEBUG_moreThan20Secs = false

	levelPrototype *ebiten.Image
	levelProtoDbg  []*ebiten.Image

	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (g *Game) Init() {
	g.startTime = time.Now()
	emptyImage.Fill(color.White)
	g.levelOverlayLayers = append(g.levelOverlayLayers,
		ebiten.NewImage(32*32, 32*32))
	g.levelOverlayLayers = append(g.levelOverlayLayers,
		ebiten.NewImage(32*32, 32*32))
	g.levelOverlayLayers = append(g.levelOverlayLayers,
		ebiten.NewImage(32*32, 32*32))
	g.levelOverlayLayers = append(g.levelOverlayLayers,
		ebiten.NewImage(32*32, 32*32))
	g.debugOverlay =
		ebiten.NewImage(32*32, 32*32)
	levelPrototype, levelProtoDbg = GenerateLevel(g)

}

func (g *Game) Update() error {
	var newBalls []Ball
	for _, ball := range g.Balls {
		b := BallMoveTick(ball, g)
		if !b.CanBeRemoved {
			newBalls = append(newBalls, b)
		}
	}
	g.Balls = newBalls
	mx, my := ebiten.CursorPosition()

	px := g.Player.X //- float64(g.lastFrameSx) //- sx
	py := g.Player.Y // - float64(g.lastFrameSy) //- sy
	//px *= 2
	//py *= 2
	g.debugOverlay.Clear()

	ebitenutil.DrawLine(g.debugOverlay, px, py, float64(mx), float64(my), color.RGBA{255, 0, 0, 255})
	//time for math:tm!
	dx := math.Abs(float64(mx) - px)
	dy := math.Abs(float64(my) - py)
	netDist := math.Sqrt((dx * dx) + (dy * dy))

	//ebitenutil.DebugPrintAt(g.levelOverlayLayers[4], fmt.Sprintf("%v | %v", g.lastFrameSx, g.lastFrameSy), mx, my)
	accuracyRadius := (20 * (netDist / 100))

	//biased random point (sqrt random would make it uniform)
	accOffset := accuracyRadius * rand.Float64()
	accTheta := rand.Float64() * 2 * math.Pi
	targetX := float64(mx) + (accOffset * math.Cos(accTheta))
	targetY := float64(my) + (accOffset * math.Sin(accTheta))

	ebitenutil.DrawLine(g.debugOverlay, float64(mx), float64(my), targetX, targetY, color.RGBA{255, 0, 255, 255})
	ebitenutil.DrawLine(g.debugOverlay, px, py, targetX, targetY, color.RGBA{0, 0, 255, 255})

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Balls = append(g.Balls, Ball{
			X:     px,
			Y:     py,
			DestX: targetX,
			DestY: targetY,
			Speed: 5,
		})
	}
	//clear
	if len(g.Balls) > 1600 {
		g.Balls = g.Balls[5:]
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		levelPrototype, levelProtoDbg = GenerateLevel(g)
	}

	for _, key := range inpututil.PressedKeys() {
		switch key {
		case ebiten.KeyW:
			//g.Player.Y -= 7
			g.Player.momY -= 2
		case ebiten.KeyS:
			g.Player.momY += 2

			//g.Player.Y += 7
		case ebiten.KeyA:

			g.Player.momX -= 2
			//g.Player.X -= 7
		case ebiten.KeyD:

			g.Player.momX += 2
			//g.Player.X += 7
		}
	}

	gLX := g.Player.X
	gLY := g.Player.Y

	g.Player.X += g.Player.momX
	g.Player.Y += g.Player.momY
	g.Player.X = math.Floor(g.Player.X)
	g.Player.Y = math.Floor(g.Player.Y)
	if g.Player.momX > 10 {
		g.Player.momX = 10
	}
	if g.Player.momX < -10 {
		g.Player.momX = -10
	}
	if g.Player.momY > 10 {
		g.Player.momY = 10
	}
	if g.Player.momY < -10 {
		g.Player.momY = -10
	}
	g.Player.momX -= 0.05 * g.Player.momX
	g.Player.momY -= 0.05 * g.Player.momY

	pTX := int(g.Player.X / 32)
	pTY := int(g.Player.Y / 32)

	if g.Level[pTX][pTY].T == 10 {
		g.Player.momX = 0
		g.Player.momY = 0
		g.Player.X = gLX
		g.Player.Y = gLY
	}

	//TODO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//Add A HITBOX system for the love of all that is good

	if g.Player.X < 0 {
		g.Player.momX = 0
		g.Player.X = 0
	}
	if g.Player.Y < 0 {
		g.Player.momY = 0
		g.Player.Y = 0
	}
	if g.Player.X > 960-16 {
		g.Player.momX = 0
		g.Player.X = 960 - 16
	}
	if g.Player.Y > 720-16 {
		g.Player.momY = 0
		g.Player.Y = 720 - 16
	}

	if time.Now().Sub(g.startTime).Seconds() > 20 {
		DEBUG_moreThan20Secs = true
	}

	return nil
}

func DrawLightCircle(dst *ebiten.Image, cx, cy, r float64, clr color.Color) {
	var path vector.Path
	rd, g, b, a := clr.RGBA()

	path.Arc(float32(cx), float32(cy), float32(r), 0, 2*math.Pi, vector.Clockwise)

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = float32(rd) / 0xffff
		vertices[i].ColorG = float32(g) / 0xffff
		vertices[i].ColorB = float32(b) / 0xffff
		vertices[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{}
	op.CompositeMode = ebiten.CompositeModeLighter

	dst.DrawTriangles(vertices, indices, emptySubImage, nil)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 255})
	g.levelOverlayLayers[2].Fill(color.RGBA{15, 15, 15, 255})

	DrawLightCircle(g.levelOverlayLayers[2], g.Player.X+16, g.Player.Y+16, 128, color.RGBA{255, 245, 245, 100})
	//DrawLightCircle(g.levelOverlayLayers[2], g.Player.X+16, g.Player.Y+16, 22, color.RGBA{255, 255, 255, 150})

	for i, layer := range g.levelOverlayLayers {
		if i == 2 {
			op := &ebiten.DrawImageOptions{}
			op.CompositeMode = ebiten.CompositeModeMultiply
			screen.DrawImage(layer, op)
		} else {
			screen.DrawImage(layer, nil)
		}
	}
	//terrain
	for x, r := range g.Level {
		for y, t := range r {
			if t.IsSalted {
				ebitenutil.DrawRect(screen, float64(x)*32, float64(y)*32, 32, 32, color.RGBA{0, 255, 0, 100})
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))
	screen.DrawImage(g.Assets.Img["player"], op)

	//ebitenutil.DrawLine(screen, float64(g.Player.X), float64(g.Player.Y), float64(mx), float64(my), color.RGBA{255, 0, 0, 255})

	for _, ball := range g.Balls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(ball.X, ball.Y)

		screen.DrawImage(g.Assets.Img["golfball"], op)
		//ebitenutil.DrawCircle(screen, ball.X, ball.Y, 2, color.RGBA{255, 10, 255, 255})

	}

	//now, for a magic trick!
	//scrolling: minimal pain (hopefully) edition!

	AP_X_SZ := 960.
	AP_Y_SZ := 720.

	sx := g.Player.X - (AP_X_SZ / 4)
	sy := g.Player.Y - (AP_Y_SZ / 4)
	ex := g.Player.X + (AP_X_SZ / 4)
	ey := g.Player.Y + (AP_Y_SZ / 4)

	if sx < 0 {
		sx = 0
		ex = AP_X_SZ / 2
	}
	if sy < 0 {
		sy = 0
		ey = AP_Y_SZ / 2
	}
	if ex > 960 {
		ex = AP_X_SZ
		sx = AP_X_SZ - (AP_X_SZ / 2)
	}
	if ey > 720 {
		ey = AP_Y_SZ
		sy = AP_Y_SZ - (AP_Y_SZ / 2)
	}

	g.lastFrameSx = int(sx)
	g.lastFrameSy = int(sy)

	subview := ebiten.NewImageFromImage(screen.SubImage(image.Rect(int(sx), int(sy), int(ex), int(ey))))
	op = &ebiten.DrawImageOptions{}

	op.GeoM.Scale(960/AP_X_SZ*2, 720/AP_Y_SZ*2)
	screen.Clear()
	//ebitenutil.DrawRect(screen, 20, 20, 600, 400, color.Black)
	screen.DrawImage(subview, op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("netX %v netY %v (> %v - %v | ^ %v - %v)", ex-sx, ey-sy, sx, ex, sy, ey), 120, 100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %v (ball count %v) | TPS: %v", int(ebiten.CurrentFPS()), len(g.Balls), int(ebiten.CurrentTPS())), 430, 10)

	//== UI ==//

	//ebitenutil.DrawRect(screen, float64(g.Player.X), float64(g.Player.Y), 16, 16, color.RGBA{200, 10, 200, 255})
	mx, my := ebiten.CursorPosition()

	px := g.Player.X - sx //- sx
	py := g.Player.Y - sy //- sy
	px *= 2
	py *= 2
	//time for math:tm!
	dx := math.Abs(float64(mx) - px)
	dy := math.Abs(float64(my) - py)
	netDist := math.Sqrt((dx * dx) + (dy * dy))
	radius := (20 * (netDist / 100))
	//debug print it
	//ebitenutil.DrawCircle(screen, float64(mx), float64(my), radius, color.RGBA{0, 255, 255, 255})
	radius /= 8
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	op.GeoM.Translate(float64(mx)-8*radius, float64(my)-8*radius)
	ebitenutil.DrawLine(screen, px, py, float64(mx), float64(my), color.RGBA{0, 255, 0, 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("circle @ (%v,%v) | Mouse @ (%v,%v)\bPlayer @ (%v,%v)", float64(mx)-8*radius, float64(my)-8*radius, mx, my, px, py), mx, my)
	if radius < 3.2 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm2"], op)
	} else if radius < 8 {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm3"], op)
	} else {
		screen.DrawImage(g.Assets.Img["ui/aimmarker_sm1"], op)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 10)
	screen.DrawImage(g.Assets.Img["ui/bars/time"], op)
	if time.Now().Sub(g.startTime).Seconds() <= 20 {
		ebitenutil.DrawRect(screen, 48, 22, 80-(time.Now().Sub(g.startTime).Seconds()*4), 8, color.RGBA{76, 76, 129, 255})
	} else {
		screen.DrawImage(g.Assets.Img["ui/bars/timeover"], op)
	}
	op.GeoM.Translate(0, 37)
	screen.DrawImage(g.Assets.Img["ui/bars/score"], op)
	//todo: actual font
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %v (%v)", g.Score, float64(g.Score)/float64(g.BestPossibleScore)), 40, 54)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(0, 480-32)
	ebitenutil.DebugPrintAt(screen, "LEVEL LAYOUT:", 0, 420)
	screen.DrawImage(levelPrototype, op)
	screen.DrawImage(g.debugOverlay, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 960, 720
}
