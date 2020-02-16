# texiot

`texiot` is my submission for the ETHDenver 2020 buildathon. It is an IoT mesh network comprised of multiple ESP32 chips, including an ESP32-CAM chip which streams a webcam on the meshnet. The meshnet is bridged to a MKR1000 board functioning as a regular WiFi network, giving the ability to route traffic from the meshnet to hosts on the regular WiFi network, although I didnt get a chance to implement the routing of traffic, however the meshnet is bridged to the regular wifi network. Additionally there is a backend comprised of TemporalX, S3X, and a generalized LibP2P node. TemporalX + S3X are used to pull the mjpeg stream coming off the ESP32-CAM and does a few things:

* Copies the output into a LibP2P protocol
  * This means that any libp2p host can stream the video coming off the camera
  * LibP2P video streaming, wat? :O
* Copies the output into S3X
  * The entire video feed is stored on S3 + IPFS as a backup
* Takes the data from S3X and streams it via a http server

## Accomplishes

* Lays the framework for a LibP2P IoT network
* Lays the framework for streaming video realtime over LibP2P
* Lays the framework for streaming video near-realtime over IPFS + HTTP

## Wanted To Do But No Time

* Relay the video feed over LoRa using [libp2p-lora-transport](https://github.com/RTradeLtd/libp2p-lora-transport)
* Enable routing of traffic to/from the meshnet and the regular wifi network
* Bridge the LibP2P network from the meshnet to the public network

## Related Source Files

* `src/esp32/mesh/bridge.ino`
  * The esp-mesh bridge
* `src/esp32/mesh/camera.ino`
  * The esp-mesh esp32-cam node
* `src/mkr1000/router.ino`
  * The MKR1000 wifi router, acting as the "normal" network
* `src/golang/fetcher.go`
  * The backend service
  * Reads data from the esp32-cam node and feeds into:
    * LibP2P
    * S3 + IPFS
    * HTTP Server

# Architecture

## Hardware

* Arduino MKR1000 is a normal WiFi access point
* 1x ESP32 acts as the esp-mesh root node, and a bridge
  * Connects to other ESPs in a mesh
  * Bridges normal WiFi (aka, the internet) and the mesh network
* 1x ESP32-CAM acts an an esp-mesh node, and a streaming web cam

## Diagram

```
----------------       ------------         --------------
| esp32-cam *1 | <---> | esp32 *2 | <-----> | MKR1000 *3 |
----------------       ------------         --------------
      ^                     ^
      |                     |
      -----------------------
                ^
                |
                v           
            -------------
            | laptop *4 |
            -------------


Legend:

1) ESP-MESH node
2) ESP-MESH root node & bridge
3) MKR1000 normal WiFi (meant to emulate another network, IoT meshnet, internet, etc..)
4) connects to either of the ESP-MESH nodes, runs TemporalX, S3X, LibP2P streamer, HTTP streamer
```


# Setup

* Power on MKR1000
  * Verify serial works
* Power on ESP-MESH root node / bridge
  * Verify serial works
* Power on ESP-MESH cam node
  * Verify serial works
* Connect laptop to meshnet
* Start `fetcher`

## Issues

### esp32-cam

Due to the amount of work being performed by the esp32-cam chipset which has 1 core, vs 2, and a smaller board, it overheats and crashes every 4-5 minutes requiring it be restarted. 

### general

due to the amount of traffic on the meshnet, sometimes the entire mesh seems to crash or hang. Generally this requires a hard restart, which means unplugging everything and setting it up from scratch. Occassionally a soft restart works (aka resetarting all the devices with the `RST` button)