package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

func drawDick(screen tcell.Screen, pos_x int) {
	var theDick = `
				                    			  .|||||.
                             |||||||||
                            ||||||  .
                            ||||||   >
                           ||||||| -/
                           ||||||'(
                           .'      \\
                        .-'    | | |
                       /        \ \ \
                      |       ---:.\ \
        ____________._>           \_\____ ,--.__
   ,--""           /      - .     |)_)    '\     '\
  /  "             |      .-'     /          \      '\
,/                  \           .'            '\     |
| "   "   "          \         /                '\,  /
|            " , =_____-.   .-'_________________,--"
|  "    "    /"/'      /\>-' ( <
\  "      ",/ /       ( <    |\_)
 \   ",",_/,-'        |\_)
  '-;.__:-'`

	style := tcell.StyleDefault
	rgb := tcell.NewHexColor(int32(rand.Int() & 0xffffff))
	style = style.Foreground(rgb)

	lines := strings.Split(theDick, "\n")

	for y, line := range lines {
		for x, cell := range line {
			screen.SetContent(pos_x+x, y, cell, []rune(""), style)
		}
	}
	screen.Show()
	screen.Clear()
}

func main() {

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.Clear()

	pos_x := 0
	screen_width, _ := screen.Size()

	quit := make(chan struct{})
	go func() {
		for {
			event := screen.PollEvent()
			switch event := event.(type) {
			case *tcell.EventKey:
				switch event.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					close(quit)
					return
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 15):
		}
		if screen_width == pos_x {
			break loop
		}
		drawDick(screen, pos_x)
		pos_x++
	}
	screen.Fini()
	command_to_run := strings.Join(os.Args[1:], " ")
	fmt.Printf("Running docker comand: docker %s\n", command_to_run)
	cmd := exec.Command("docker", command_to_run)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
