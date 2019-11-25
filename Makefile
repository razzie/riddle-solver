.PHONY: riddle-solver
.DEFAULT_GOAL := riddle-solver

riddle-solver:
	go build -ldflags="-s -w" -gcflags=-trimpath=$(CURDIR) ./cmd/riddle-solver
