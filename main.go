package main

import (
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

type Dire int

const (
	Up Dire = iota
	Down
	Right
	Left
)

type Pixel struct {
	BgColor termbox.Attribute
}

type Player struct {
	X         int
	Y         int
	Direction Dire
}

var screen [][]Pixel
var maxC, maxL int
var snake [][]int
var feeds [][]int
var character Player
var running = true
var isSnakeFeeding = false
var score = 0
var lostMessage = "You lost! Goodbye!"

func initialize() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	maxC, maxL = termbox.Size()

	screen = make([][]Pixel, maxL)
	for i := range screen {
		screen[i] = make([]Pixel, maxC)
		for j := range screen[i] {
			screen[i][j] = Pixel{
				BgColor: termbox.ColorBlue,
			}
		}
	}
	for i := 0; i < 5; i++ {
		feeds = append(feeds, make([]int, 5))
		feeds[i] = []int{rand.IntN(maxL), rand.IntN(maxC)}
	}
	return nil
}

func drawLose() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i, ch := range lostMessage {
		termbox.SetCell(int(maxC/2)+i-int(len(lostMessage)/2), int(maxL/2)-1, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
	time.Sleep(3 * time.Second)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func checkLose() {
	for i := range len(snake) - 2 {
		if character.Y == snake[i][0] && character.X == snake[i][1] {
			running = false
		}
	}
}

func drawPlayer() {
	if character.Direction == Up {
		character.Y--
	}
	if character.Direction == Down {
		character.Y++
	}
	if character.Direction == Right {
		character.X++
	}
	if character.Direction == Left {
		character.X--
	}
	if character.Y < 0 {
		character.Y = maxL - 1
	}
	if character.X < 0 {
		character.X = maxC - 1
	}
	character.Y = character.Y % maxL
	character.X = character.X % maxC
	snake = append(snake, []int{character.Y, character.X})
	if isSnakeFeeding && score == 1 {
		snake = append(snake, []int{character.Y, character.X})
	}
	if isSnakeFeeding {
		isSnakeFeeding = false
	} else {
		snake = snake[1:]
	}
	for i := range snake {
		screen[snake[i][0]][snake[i][1]] = Pixel{BgColor: termbox.ColorLightBlue}
	}
	screen[character.Y][character.X] = Pixel{BgColor: termbox.ColorBlue}
}

func drawFeed() {
	for i := 0; i < len(feeds); i++ {
		if feeds[i][0] == character.Y && feeds[i][1] == character.X {
			randoms := []int{rand.IntN(maxL), rand.IntN(maxC)}
			for i := range snake {
				if snake[i][0] == randoms[0] && snake[i][1] == randoms[1] {
					randoms = []int{rand.IntN(maxL), rand.IntN(maxC)}
				}
			}
			feeds[i] = randoms
			score++
			isSnakeFeeding = true
		}
		screen[feeds[i][0]][feeds[i][1]] = Pixel{BgColor: termbox.ColorRed}
	}
}

func render() {
	for i := range screen {
		for j := range screen[i] {
			screen[i][j] = Pixel{
				BgColor: termbox.ColorGreen,
			}
		}
	}
	drawFeed()
	drawPlayer()
	checkLose()
	for y := 0; y < maxL && y < len(screen); y++ {
		for x := 0; x < maxC && x < len(screen[y]); x++ {
			pixel := screen[y][x]
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, pixel.BgColor)
		}
	}
	scoreStr := strconv.Itoa(score)

	for i, ch := range scoreStr {
		termbox.SetCell(0+i, 0, ch, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.Flush()
}

func main() {
	err := initialize()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	character = Player{X: 0, Y: 0, Direction: Right}

	var event termbox.Event
	go func() {
		for running {
			event = termbox.PollEvent()
			if event.Type == termbox.EventKey && event.Key == termbox.KeyEsc {
				running = false
				lostMessage = "Goodbye!"
				return
			}
			if event.Type == termbox.EventKey && event.Key == termbox.KeyArrowUp && character.Direction != Down {
				character.Direction = Up
			}
			if event.Type == termbox.EventKey && event.Key == termbox.KeyArrowDown && character.Direction != Up {
				character.Direction = Down
			}
			if event.Type == termbox.EventKey && event.Key == termbox.KeyArrowLeft && character.Direction != Right {
				character.Direction = Left
			}
			if event.Type == termbox.EventKey && event.Key == termbox.KeyArrowRight && character.Direction != Left {
				character.Direction = Right
			}
		}
	}()

	for running {
		time.Sleep(200 * time.Millisecond)
		render()
	}
	drawLose()
}
