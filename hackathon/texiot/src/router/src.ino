#include <WiFi.h>
#include <WebServer.h>
#include <ESPmDNS.h>
#include <esp_system.h>
#include <esp_timer.h>

#define AP_SSID "temporalxisbest"
#define AP_PASS "temporalxisbest123"
#define AP_CHAN 14
#define AP_HIDDEN 0
#define AP_MAX_CONNS 12

IPAddress ip(192, 168, 0, 1);
IPAddress gateway(192, 168, 0, 1);
IPAddress subnet(255, 255, 255, 0);

// create server that listens on port 80
WebServer server(80);

/*
 * https://github.com/espressif/arduino-esp32/blob/master/libraries/WiFi/examples/WiFiClientEvents/WiFiClientEvents.ino#L144
 * https://github.com/espressif/arduino-esp32/blob/master/libraries/WiFi/src/WiFiAP.h
 * 
 */

void handleDebug() {
  String message = "";
  message += "reason of last reset: ";
  message += esp_reset_reason();
  message += "\n";
  message += "available heap size: ";
  message += esp_get_free_heap_size();
  message += "\n";
  server.send(200, "text/plain", message);
}

void handleNotFound() {
  String message = "File Not Found\n\n";
  message += "URI: ";
  message += server.uri();
  message += "\nMethod: ";
  message += (server.method() == HTTP_GET) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server.args();
  message += "\n";
  for (uint8_t i = 0; i < server.args(); i++) {
    message += " " + server.argName(i) + ": " + server.arg(i) + "\n";
  }
  server.send(404, "text/plain", message);
}

void handleEvent(WiFiEvent_t event) {
  switch (event) {
    case SYSTEM_EVENT_AP_STACONNECTED:
      Serial.println("client connected");
      break;
    case SYSTEM_EVENT_AP_STADISCONNECTED:
      Serial.println("client disconnected");
      break;
    case SYSTEM_EVENT_AP_STAIPASSIGNED:
      Serial.println("assignedip to client");
      break;
    case SYSTEM_EVENT_AP_PROBEREQRECVED:
      Serial.println("received probe request");
      break;
  }
}

void setup() {
  Serial.begin(115200);
  while (!Serial) ;
  // enable access point + station mode
  WiFi.mode(WIFI_MODE_APSTA);
  // enable long range comms
  WiFi.enableLongRange(true);
  // set tx power to max
  WiFi.setTxPower(WIFI_POWER_19_5dBm);
  // register call back handler
  if (!WiFi.onEvent(handleEvent)) {
    Serial.println("failed to set event handler");
    while (true) ;
  }
  // configure the access point
  if (!WiFi.softAPConfig(ip, gateway, subnet)) {
    Serial.println("failed to set soft ap config");
    while (true) ;
  }
  // enable the access point
  if (!WiFi.softAP(AP_SSID, AP_PASS, AP_CHAN /* channel */, AP_HIDDEN /* hidden */, AP_MAX_CONNS /* max conns*/)) {
    Serial.println("failed to initialize access point");
    while (true) ;
  }
  // register mdns
  if (!MDNS.begin("esp32-texiot")) {
    Serial.println("failed to start mdns responder");
    while (true) ;
  }
  // register web server handler and start
  server.on("/", []() {
    server.send(200, "text/plain", "temporalx iot");
  });
  server.on("/debug", []() {
    handleDebug();
  });
  server.onNotFound(handleNotFound);
  server.begin();
  Serial.println("access point initialized");
  Serial.flush();
}

void loop() {
  // if no client is available this doesn't block
  server.handleClient();
}
