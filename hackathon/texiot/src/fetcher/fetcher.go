package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// taken from: https://gist.github.com/ZenGround0/af448f56882c16aaf10f39db86b4991e

func main() {
	tr := http.DefaultTransport
	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	r := os.Stdin
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9094",
			Path:   "/",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          r,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024*1024) // 1mb buffer
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			log.Println("received eof from server, stream ended")
			os.Exit(0)
		}
		fmt.Printf("read %v bytes", n)
	}
}
