# texiot

`texiot` is my submission for the ETHDenver 2020 buildathon. It is an IoT mesh network comprised of multiple ESP32 chips, including an ESP32-CAM chip which streams a webcam on the meshnet. Additionally there is a backend comprised of TemporalX, S3X, and a generalized LibP2P node. TemporalX + S3X are used to pull the mjpeg stream coming off the ESP32-CAM and does a few things:

* Copies the output into a LibP2P protocol
  * This means that any libp2p host can stream the video coming off the camera
  * LibP2P video streaming, wat? :O
* Copies the output into S3X
  * The entire video feed is stored on S3 + IPFS as a backup
* Takes the data from S3X and streams it via a http server

This submission essentially accomplishes a few different things:

* Lays the framework for a LibP2P IoT network
* Lays the framework for streaming video realtime over LibP2P
* Lays the framework for streaming video near-realtime over IPFS + HTTP

# Architecture

## Hardware

* Arduino MKR1000 is a normal WiFi access point
* 1x ESP32 acts as the esp-mesh root node, and a bridge
  * Connects to other ESPs in a mesh
  * Bridges normal WiFi (aka, the internet) and the mesh network
* 1x ESP32 acts as an esp-mesh node
* 1x ESP32-CAM acts an an esp-mesh node, and a streaming web cam

## Software

* "Fetcher" program fetches data from the streaming web cam and:
  * 1. Sends it over a LibP2P protocol
  * 2. Stores an archived feed of the video on IPFS via S3

# Demo

To demo this manually you need to downlaod the video feed object from s3 or ipfs and save as an `.mjpeg` object. Then you must use `ffmpeg` to convert to `mp4` and the contents can be viewed.

Example command:

```shell
$> ffmpeg -i videofeed.mjpeg -c:v libx264 -preset veryslow -crf 18 output.mp4
```

Ideally we can use a live encoder to fetch the data via the libp2p protocol, encode to mp4 live and display.

# Ideas

## Stream Encoder

* Retrieve data from libp2p
* Encode
* Serve via http

Links:

* https://stackoverflow.com/questions/60168799/saving-a-continuous-stream-of-images-from-ffmpeg-image2pipe
* https://github.com/mattn/go-mjpeg
* https://github.com/bonedaddy/mjpeg
* https://code.nfsmith.ca/nsmith/mjpeg
* https://github.com/gen2brain/cam2ip
* https://github.com/as/video
* https://github.com/panzerdev/mjpeg-stitcher

Notes:

Possible `io.Copy` implementation of retrieval + encode

```Go
// from https://stackoverflow.com/questions/43601846/golang-and-ffmpeg-realtime-streaming-input-output

var buf bytes.Buffer

n, err := io.Copy(&buf, stdout)
verificaErro(err)
fmt.Printf("Copied %d bytes\n", n)

err = cmd.Wait()
fmt.Printf("Wait error %v\n", err)

//do something with the data
data := buf.Bytes()
f, err4 := os.OpenFile(dir+"/out.raw", os.O_RDWR|os.O_APPEND, 0666)
verificaErro(err4)
defer f.Close()
nw, err := f.Write(data)
f.Sync()
fmt.Printf("Write size %d bytes\n", nw)
```

Using [`thumbnailer`](https://github.com/bakape/thumbnailer)

```Go
// from https://stackoverflow.com/questions/49800771/piping-raw-byte-video-to-ffmpeg-go
thumbnailDimensions := thumbnailer.Dims{Width: 320, Height: 130}

thumbnailOptions := thumbnailer.Options{JPEGQuality:100, MaxSourceDims:thumbnailer.Dims{}, ThumbDims:thumbnailDimensions, AcceptedMimeTypes: nil}

sourceData, thumbnail, err := thumbnailer.ProcessBuffer(videoData, thumbnailOptions)

imageBytes := thumbnail.Image.Data
```