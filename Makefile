.PHONY: riddle-solver
.DEFAULT_GOAL := riddle-solver

riddle-solver:
	go build -ldflags="-s -w" ./cmd/riddle-solver
