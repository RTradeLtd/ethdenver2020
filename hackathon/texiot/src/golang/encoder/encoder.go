package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mattn/go-mjpeg"
)

var (
	fileName = "videofeed.mjpeg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Fatal(err)
	}
	fh.Close()
	stream := mjpeg.NewStreamWithInterval(time.Second)
	if err := stream.Update(append(data[0:0:0], data...)); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if err := stream.Update(append(data[0:0:0], data...)); err != nil {
				log.Fatal(err)
			}
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/mjpeg", stream.ServeHTTP)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<img src="/mjpeg" />`))
	})
	server := &http.Server{Addr: "0.0.0.0:6969", Handler: mux}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		cancel()
		server.Shutdown(ctx)
	}()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	stream.Close()
}
