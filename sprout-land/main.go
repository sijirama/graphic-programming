package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 900 //1000
	ScreenHeigth = 680 //480
)

var (
	gameStateRunning = true
	gameBgColor      = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	hillSprite   rl.Texture2D
	fenceSprite  rl.Texture2D
	houseSprite  rl.Texture2D
	waterSprite  rl.Texture2D
	tilledSprite rl.Texture2D
	tex          rl.Texture2D

	playerSprite rl.Texture2D

	playerSrc  rl.Rectangle
	playerDest rl.Rectangle

	playerSpeed      float32 = 1.3
	playerSpeedBoost float32 = 2

	playerMoving bool
	playerDir    int
	playerFrame  int

	playerUp, playerDown, playerLeft, playerRight bool

	frameCount int

	tileDest   rl.Rectangle // draw the tile on the screen
	tileSrc    rl.Rectangle // src of the image
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	musicPaused bool
	music       rl.Music

	camera rl.Camera2D
)

func DrawScene() {
	//rl.DrawTexture(grassSprite, 100, 50, rl.White)

	for i := 0; i < len(tileMap); i++ {
		if (tileMap[i]) != 0 {
			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapW)

			if srcMap[i] == "g" {
				tex = grassSprite
			}
			if srcMap[i] == "l" {
				tex = hillSprite
			}
			if srcMap[i] == "f" {
				tex = fenceSprite
			}
			if srcMap[i] == "h" {
				tex = houseSprite
			}
			if srcMap[i] == "w" {
				tex = waterSprite
			}
			if srcMap[i] == "t" {
				tex = tilledSprite
			}

			if srcMap[i] == "h" || srcMap[i] == "f" {
				tileSrc.X = 0
				tileSrc.Y = 0
				rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0.0, rl.White)
			}

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))

			rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0.0, rl.White)
		}
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0.0, rl.White)
}

func Input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}
	if rl.IsKeyDown(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func Update() {
	gameStateRunning = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}
	camera.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerLeft, playerRight = false, false, false, false
}

func Render() {
	rl.BeginDrawing()
	rl.ClearBackground(gameBgColor)
	rl.BeginMode2D(camera)

	DrawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile)

	if err != nil {
		fmt.Printf("Failed to open map file: %s", err)
		os.Exit(1)
	}

	remNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")

	mapW, mapH = -1, -1

	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < (mapW*mapH)+2 {
			tileMap = append(tileMap, m*56)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}

	//remove extra error inputs from file
	if len(tileMap) > (mapW * mapH) {
		tileMap = tileMap[:len(tileMap)-1]
	}
}

func init() {
	rl.InitWindow(ScreenWidth, ScreenHeigth, "RayRay")
	rl.InitAudioDevice()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	camera = rl.NewCamera2D(rl.NewVector2(float32(ScreenWidth/2), float32(ScreenHeigth/2)), rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)), float32(playerDest.Y-(playerDest.Height/2))), 0.0, 2.5)

	camera.Zoom = 3

	// Give audio system time to initialize
	time.Sleep(100 * time.Millisecond)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	fenceSprite = rl.LoadTexture("assets/Tilesets/Fences.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/Wooden_House.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled_Dirt.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("assets/Characters/BasicCharacterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 60, 60)

	music = rl.LoadMusicStream("assets/song.mp3")
	musicPaused = true
	rl.PlayMusicStream(music)

	loadMap("two.map")
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {
	for gameStateRunning {
		Input()
		Update()
		Render()
	}

	quit()
}
