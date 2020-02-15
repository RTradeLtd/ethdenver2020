ifeq ($(BOARD),)
BOARD := arduino:samd:mkr1000
endif

ifeq ($(PORT),)
PORT := /dev/ttyACM0
endif

.PHONY: build
build: compile-sketch upload-sketch

.PHONY: compile-sketch
compile-sketch:
	arduino-cli compile --fqbn  $(BOARD) src

.PHONY: upload-sketch
upload-sketch:
	arduino-cli upload -p $(PORT) --fqbn $(BOARD) src

.PHONY: list-boards
list-boards:
	arduino-cli board list

.PHONY: update-index
update-index:
	arduino-cli core update-index

.PHONY: install-libs
install-libs:
	arduino-cli lib install ArduinoUni
