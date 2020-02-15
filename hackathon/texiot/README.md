# texiot

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

