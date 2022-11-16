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
)

const (
	USE_TIMECHECK = true
	DEBUG_ROTVIEW = false
)

var (
	DEBUG_moreThan20Secs  = false
	DEBUG_showLevelLayers = false
	DEBUG_radians         = 0.
	levelPrototype        *ebiten.Image

	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (g *Game) Init() {
	g.startTime = time.Now()
	emptyImage.Fill(color.White)
	for i := 0; i < 4; i++ {
		g.levelOverlayLayers = append(g.levelOverlayLayers,
			ebiten.NewImage(32*32, 32*32))
	}
	g.colOverlay = ebiten.NewImage(32*32, 32*32)
	g.debugOverlay = ebiten.NewImage(32*32, 32*32)
	levelPrototype = GenerateLevel(g)
	if DEBUG_ROTVIEW {
		ebiten.SetWindowTitle("SaltRush ~ Rot View Radian Rotation Debugger")
	}
}

func (g *Game) Update() error {
	g.debugOverlay.Clear()
	g.colOverlay.Clear()

	var newBalls []Ball
	for _, ball := range g.Balls {
		b := BallMoveTick(ball, g)
		if !b.CanBeRemoved {
			newBalls = append(newBalls, b)
		}
	}
	g.Balls = newBalls
	mx, my := ebiten.CursorPosition()

	mx /= 2
	my /= 2
	mx += g.lastFrameSx
	my += g.lastFrameSy

	px := g.Player.X
	py := g.Player.Y

	//time for math:tm!
	dx := math.Abs(float64(mx) - px)
	dy := math.Abs(float64(my) - py)
	netDist := math.Sqrt((dx * dx) + (dy * dy))

	accuracyRadius := (20 * (netDist / 100))

	accOffset := accuracyRadius * rand.Float64()
	accTheta := rand.Float64() * 2 * math.Pi
	targetX := float64(mx) + (accOffset * math.Cos(accTheta))
	targetY := float64(my) + (accOffset * math.Sin(accTheta))

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
		levelPrototype = GenerateLevel(g)
	}
	propPx := g.Player.X
	propPy := g.Player.Y
	for _, key := range inpututil.PressedKeys() {
		switch key {
		case ebiten.KeyW:
			propPy -= 5

			//g.Player.momY -= 2
		case ebiten.KeyS:
			propPy += 5

			//g.Player.momY += 2
		case ebiten.KeyA:
			g.Player.Rotation -= 0.1
			//propPx -= 5
			//g.Player.Dir = "l"
			//g.Player.momX -= 2
		case ebiten.KeyD:
			g.Player.Rotation += 0.1
			//propPx += 5
			//g.Player.Dir = "r"
			//g.Player.momX += 2
		case ebiten.KeyL:
			DEBUG_showLevelLayers = !DEBUG_showLevelLayers
		case ebiten.Key1:
			DEBUG_radians -= 0.005 * math.Pi
			if DEBUG_radians < -2*math.Pi {
				DEBUG_radians = 0
			}
		case ebiten.Key2:
			DEBUG_radians += 0.005 * math.Pi
			if DEBUG_radians > 2*math.Pi {
				DEBUG_radians = 0
			}
		}
	}
	spRot := math.Abs(math.Mod(g.Player.Rotation, math.Pi))
	ebitenutil.DebugPrintAt(g.colOverlay, fmt.Sprintf("R = %v [ MOD %v] ", g.Player.Rotation, spRot), int(g.Player.X), int(g.Player.Y))
	if spRot > 5.4 || spRot < 1 {
		g.Player.Dir = "u"
	} else if spRot >= 4.3 && spRot <= 5.4 {
		g.Player.Dir = "l"
	} else if spRot <= 4.3 && spRot >= 2 {
		g.Player.Dir = "d"
	} else {
		g.Player.Dir = "r"
	}

	// propPx := g.Player.X + g.Player.momX
	// propPy := g.Player.Y + g.Player.momY
	propPx = math.Floor(propPx)
	propPy = math.Floor(propPy)
	validMovement := true
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

	for x := int(g.Player.X/32) - 5; x < int(g.Player.X/32)+5; x++ {
		for y := int(g.Player.Y/32) - 5; y < int(g.Player.Y/32)+5; y++ {
			if x < 0 || x >= 30 || y < 0 || y >= 24 {
				continue
			}

			if g.Level[x][y].T == 10 {
				if x*32 <= int(propPx)+32 &&
					x*32+32 >= int(propPx) &&
					y*32 <= int(propPy)+32 &&
					32+y*32 >= int(propPy) {
					ebitenutil.DrawRect(g.colOverlay, float64(x)*32, float64(y)*32, 32, 32, color.RGBA{255, 0, 0, 50})
					ebitenutil.DrawRect(g.colOverlay, g.Player.X, g.Player.Y, 32, 32, color.RGBA{255, 255, 0, 50})
					validMovement = false
				}
				//	ebitenutil.DrawRect(g.colOverlay, float64(x)*32, float64(y)*32, 32, 32, color.RGBA{255, 0, 0, 50})
			}
		}
	}

	if validMovement {
		g.Player.X = propPx
		g.Player.Y = propPy
	} else {
		g.Player.momX = 0
		g.Player.momY = 0
	}

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

func (g *Game) Draw(screen *ebiten.Image) {

	//== Titlescreen (State 0)
	if g.GameState == 0 {
		screen.Fill(color.RGBA{40, 40, 40, 255})
		ebitenutil.DebugPrint(screen, "This is going to be a titlescreen")
		ebitenutil.DebugPrintAt(screen, "Press [X] to start", 0, 20)
		return //no other draw logic for title
	}
	//== Gameplay (State 1)

	if DEBUG_showLevelLayers {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.1, 0.1)
		op.GeoM.Translate(0, 32*32/10)
		for _, l := range g.levelOverlayLayers {
			g.debugOverlay.DrawImage(l, op)
			op.GeoM.Translate(0, 32*32/10)
		}
	}
	screen.Fill(color.RGBA{40, 40, 40, 255})
	g.levelOverlayLayers[3].Fill(color.RGBA{15, 15, 15, 255})

	DrawLightCircle(g.levelOverlayLayers[3], g.Player.X+16, g.Player.Y+16, 124, color.RGBA{255, 245, 245, 40})
	DrawLightCircle(g.levelOverlayLayers[3], g.Player.X+16, g.Player.Y+16, 118, color.RGBA{255, 245, 245, 40})
	DrawLightCircle(g.levelOverlayLayers[3], g.Player.X+16, g.Player.Y+16, 112, color.RGBA{255, 245, 245, 40})
	DrawLightCircle(g.levelOverlayLayers[3], g.Player.X+16, g.Player.Y+16, 106, color.RGBA{255, 245, 245, 40})

	for i, layer := range g.levelOverlayLayers {
		if i == 3 {
			op := &ebiten.DrawImageOptions{}
			op.CompositeMode = ebiten.CompositeModeMultiply
			screen.DrawImage(layer, op)
		} else {
			screen.DrawImage(layer, nil)
		}
	}

	w, h := g.Assets.Img[fmt.Sprintf("truck/truck_%v", g.Player.Dir)].Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(g.Player.Rotation)

	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))

	screen.DrawImage(g.Assets.Img[fmt.Sprintf("truck/truck_%v", g.Player.Dir)], op)

	g.levelOverlayLayers[2].Clear()
	for _, ball := range g.Balls {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(ball.X, ball.Y)
		g.levelOverlayLayers[2].DrawImage(g.Assets.Img["golfball"], op)
	}

	screen.DrawImage(g.colOverlay, nil)

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
	screen.DrawImage(subview, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %v (ball count %v) | TPS: %v", int(ebiten.CurrentFPS()), len(g.Balls), int(ebiten.CurrentTPS())), 630, 10)

	//== UI ==//
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

	radius /= 8
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(radius, radius)
	op.GeoM.Translate(float64(mx)-8*radius, float64(my)-8*radius)

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
	screen.DrawImage(g.debugOverlay, nil)

	if DEBUG_ROTVIEW {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-8, -8)
		op.GeoM.Rotate(float64(DEBUG_radians))
		op.GeoM.Scale(10, 10)
		op.GeoM.Translate(960/4, 720/2)

		screen.DrawImage(g.Assets.Img["tile/rot"], op)
		op.GeoM.Translate(960/4, -150)
		screen.DrawImage(g.Assets.Img["truck/truck_u"], op)
		op.GeoM.Translate(960/4, 0)
		screen.DrawImage(g.Assets.Img["truck/truck_d"], op)
		op.GeoM.Translate(-1*960/4, 260)
		screen.DrawImage(g.Assets.Img["truck/truck_l"], op)
		op.GeoM.Translate(960/4, 0)
		screen.DrawImage(g.Assets.Img["truck/truck_r"], op)
		x := ebiten.NewImage(100, 20)
		ebitenutil.DebugPrintAt(x, fmt.Sprintf("%v", DEBUG_radians), 0, 0)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(300, 100)
		screen.DrawImage(x, op)
	}

	//gameover is drawn *over* the main game render
	if g.GameState == 2 {
		//TODO
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 960, 720
}
