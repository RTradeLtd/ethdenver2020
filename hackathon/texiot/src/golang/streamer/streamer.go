package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/minio/minio-go"
)

const BUFSIZE = 1024 * 8

var (
	hostAddress   = flag.String("host", "192.168.0.2:80", "the host to connect to")
	addr          = flag.String("multi.addr", "/ip4/0.0.0.0/tcp/4006", "the multiaddr for libp2p host")
	accessKey     = flag.String("access.key", "minio", "minio access key")
	secretKey     = flag.String("secret.key", "minio123", "minio secret key")
	minioEndpoint = flag.String("minio.endpoint", "0.0.0.0:9000", "minio endpoint")
	setup         = flag.Bool("setup", true, "setup the testenv then exit")
	mux           sync.Mutex
	mc            *minio.Client
)

func main() {
	minioClient, err := minio.New(*minioEndpoint, *accessKey, *secretKey, false)
	if err != nil {
		log.Fatal("failed to access minio endpoint ", err)
	}
	mc = minioClient
	streamSrvMux := http.NewServeMux()
	streamSrvMux.HandleFunc("/", streamHandler)
	streamServer := http.Server{
		Addr:    "0.0.0.0:6969",
		Handler: streamSrvMux,
	}
	if err := streamServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	obj, err := mc.GetObject("testbucket", "videofeed", minio.GetObjectOptions{})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		"pipe:0",
		"-c:v",
		"libx264",
		"-preset",
		"veryslow",
		"-crf",
		"18",
		//"-s",
		//"WxH",
		"-f",
		"mjpeg",
		"pipe:1",
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	/*
			look into
			-s WxH -f mjpeg
		https://superuser.com/questions/685022/how-can-i-pipe-data-losslessly-to-and-from-ffmpeg
	*/
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	go io.Copy(stdin, obj)
	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, stdout)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if err := cmd.Wait(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	obj.Close()

	var buf = make([]byte, buffer.Len())
	n, err := buffer.Read(buf)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	fileSize := int(len(buf[:n]))
	reader := bytes.NewReader(buf[:n])
	if len(r.Header.Get("Range")) == 0 {

		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes 0-"+contentEnd+"/"+contentLength)
		w.WriteHeader(200)

		buffer := make([]byte, BUFSIZE)

		for {
			n, err := reader.Read(buffer)

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}

	} else {

		rangeParam := strings.Split(r.Header.Get("Range"), "=")[1]
		splitParams := strings.Split(rangeParam, "-")

		// response values

		contentStartValue := 0
		contentStart := strconv.Itoa(contentStartValue)
		contentEndValue := fileSize - 1
		contentEnd := strconv.Itoa(contentEndValue)
		contentSize := strconv.Itoa(fileSize)

		if len(splitParams) > 0 {
			contentStartValue, err = strconv.Atoi(splitParams[0])

			if err != nil {
				contentStartValue = 0
			}

			contentStart = strconv.Itoa(contentStartValue)
		}

		if len(splitParams) > 1 {
			contentEndValue, err = strconv.Atoi(splitParams[1])

			if err != nil {
				contentEndValue = fileSize - 1
			}

			contentEnd = strconv.Itoa(contentEndValue)
		}

		contentLength := strconv.Itoa(contentEndValue - contentStartValue + 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes "+contentStart+"-"+contentEnd+"/"+contentSize)
		w.WriteHeader(206)

		buffer := make([]byte, BUFSIZE)

		obj.Seek(int64(contentStartValue), 0)

		writeBytes := 0

		for {
			n, err := reader.Read(buffer)

			writeBytes += n

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			if writeBytes >= contentEndValue {
				data := buffer[:BUFSIZE-writeBytes+contentEndValue+1]
				w.Write(data)
				w.(http.Flusher).Flush()
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}
	}
}

func stdinfill(stdin io.WriteCloser) {
	fi, err := ioutil.ReadFile("music.ogg")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(stdin, bytes.NewReader(fi))
}

func runcommand() {

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "mp3", "pipe:1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	stdinfill(stdin)

	fo, err := os.Create("output.mp3")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(fo, stdout)

	defer fo.Close()

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
