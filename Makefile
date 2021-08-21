.PHONY: all gui tui
.DEFAULT_GOAL := all

all: gui tui

gui:
	go build -mod=vendor -ldflags="-s -w -H=windowsgui" -gcflags=-trimpath=$(CURDIR) ./cmd/riddle-solver-gui

tui:
	go build -mod=vendor -ldflags="-s -w" -gcflags=-trimpath=$(CURDIR) ./cmd/riddle-solver-tui
