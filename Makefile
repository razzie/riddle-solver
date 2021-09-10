.PHONY: all gui tui
.DEFAULT_GOAL := all
GO := go
BUILD := build -mod=vendor
GOROOT := $(shell go env GOROOT)
LDFLAGS := -s -w
GCFLAGS := -trimpath=$(CURDIR);$(GOROOT)/src

all: gui tui

ifeq ($(OS),Windows_NT)
gui: LDFLAGS += -H=windowsgui
endif
gui:
	$(GO) $(BUILD) -ldflags="$(LDFLAGS)" -gcflags=all="$(GCFLAGS)" ./cmd/riddle-solver-gui

tui:
	$(GO) $(BUILD) -ldflags="$(LDFLAGS)" -gcflags=all="$(GCFLAGS)" ./cmd/riddle-solver-tui
