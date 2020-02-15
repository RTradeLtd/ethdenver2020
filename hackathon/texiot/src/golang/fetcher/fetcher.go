package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

// from https://gist.github.com/ZenGround0/49e4a1aa126736f966a1dfdcb84abdae

const partBoundary = "123456789000000000000987654321"

const boundary = "\r\n--" + partBoundary + "\r\n"

func main() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	// Send http request chunk encoding the multipart message
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9094",
			Path:   "/",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	reader := multipart.NewReader(resp.Body, boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			log.Fatal(err)
		}
		var buf []byte
		n, err := part.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			log.Fatal("no data read")
		}
		// TODO(bonedaddy): send through libp2p
		fmt.Println(string(buf))
	}
}
