#include <SPI.h>
#include <WiFi101.h>

#define AP_SSID "temporalxisbest"
#define AP_PASS "temporalxisbest123"


IPAddress ip(192, 168, 0, 1);
IPAddress dns(192, 168, 0, 1);
IPAddress gateway(192, 168, 0, 1);
IPAddress subnet(192, 168, 0, 0);

char ssid[] = AP_SSID;
char pass[] = AP_PASS;

void setup() {
  Serial.begin(9600);
  while (!Serial) ;
  if (WiFi.status() == WL_NO_SHIELD) {
    Serial.println("WiFi shield not present");
    while (true);
  }
  WiFi.config(ip, dns, gateway, subnet);
  while (WiFi.beginAP(ssid, pass) != WL_AP_LISTENING) {
    // wait 10 seconds 
    delay(10000);
  }
  Serial.println("access point started");
  Serial.flush();
}

void loop() { }