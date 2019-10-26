.PHONY: riddle-solver
.DEFAULT_GOAL := riddle-solver

riddle-solver:
	go build ./cmd/riddle-solver
