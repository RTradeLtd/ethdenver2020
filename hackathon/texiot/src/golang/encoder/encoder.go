package main

import (
	"log"
	"os/exec"
)

var (
	fileName = "videofeed.mjpeg"
)

func main() {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		fileName,
		"-c:v",
		"libx264",
		"-preset",
		"veryslow",
		"-crf",
		"18",
		"output.mp4",
	)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
