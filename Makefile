.PHONY: riddle-solver
.DEFAULT_GOAL := riddle-solver

riddle-solver:
	go build -ldflags="-s -w -H=windowsgui" -gcflags=-trimpath=$(CURDIR) ./cmd/riddle-solver
