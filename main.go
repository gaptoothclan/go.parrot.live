package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var frames []string

// Load parrot frames
func init() {
	files, err := ioutil.ReadDir("./frames")
	if err != nil {
		fmt.Println("error reading files")
		return
	}

	for _, f := range files {
		fileByte, err := ioutil.ReadFile("./frames/" + f.Name())
		if err != nil {
			panic(err)
		}
		frames = append(frames, string(fileByte))
	}
}

// Hijack the http connection and stream the parrot
func parrotStream(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("Hijacking unsupported")
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	curFrame := 0
	maxFrame := len(frames)

	colors := [8]int{30, 31, 32, 33, 34, 35, 36, 37}
	rand.Seed(time.Now().Unix())

	for {
		time.Sleep(80 * time.Millisecond)
		conn.Write([]byte("\033c"))

		colorIndex := rand.Intn(len(colors))
		color := colors[colorIndex]

		coloredParrot := fmt.Sprintf("\033[0;%dm%s\033[0;%dm", color, frames[curFrame], color)
		_, err := conn.Write([]byte(coloredParrot))

		if err != nil {
			fmt.Println("Client closed connection")
			break
		}

		if curFrame+1 == maxFrame {
			curFrame = 0
		} else {
			curFrame = curFrame + 1
		}
	}
}

func main() {
	http.HandleFunc("/", parrotStream)
	http.ListenAndServe(":8888", nil)
}
