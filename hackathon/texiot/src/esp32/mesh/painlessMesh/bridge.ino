//************************************************************
// this is a simple example that uses the painlessMesh library to 
// connect to a node on another network. Please see the WIKI on gitlab
// for more details
// https://gitlab.com/painlessMesh/painlessMesh/wikis/bridge-between-mesh-and-another-network
//************************************************************
#include "painlessMesh.h"

#define   MESH_PREFIX     "temporalxisbestmesh"
#define   MESH_PASSWORD   "temporalxisbestmesh123"
#define   MESH_PORT       5555

#define   STATION_SSID     "temporalxisbest"
#define   STATION_PASSWORD "temporalxisbest123"
#define   STATION_PORT     5555

IPAddress myIP;

// prototypes
void receivedCallback( uint32_t from, String &msg );
painlessMesh  mesh;

void setup() {
  Serial.begin(115200);
  mesh.setDebugMsgTypes( ERROR | STARTUP | CONNECTION );  // set before init() so that you can see startup messages
  // Channel set to 6. Make sure to use the same channel for your mesh and for you other
  // network (STATION_SSID)
  mesh.init( MESH_PREFIX, MESH_PASSWORD, MESH_PORT, WIFI_AP_STA, 6 );
  // providing no station IP connects to the gateway
  mesh.stationManual(STATION_SSID, STATION_PASSWORD, STATION_PORT);
  // Bridge node, should (in most cases) be a root node. See [the wiki](https://gitlab.com/painlessMesh/painlessMesh/wikis/Possible-challenges-in-mesh-formation) for some background
  mesh.setRoot(true);
  // This node and all other nodes should ideally know the mesh contains a root, so call this on all nodes
  mesh.setContainsRoot(true);
  // set receive message callback
  mesh.onReceive(&receivedCallback);
}

void loop() {
  mesh.update();
  if (myIP != getLocalIP()) {
      myIP = getLocalIP();
      Serial.println("myip is: " myIP.toString())
  }
}

void receivedCallback( uint32_t from, String &msg ) {
  Serial.printf("bridge: Received from %u msg=%s\n", from, msg.c_str());
}

IPAddress getlocalIP() {
  return IPAddress(mesh.getStationIP());
}