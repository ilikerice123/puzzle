package game

import "fmt"

// PuzzlePoolBase represents a PuzzlePool interface
type PuzzlePoolBase interface {
	AddPuzzle(LivePuzzleBase)

	GetPuzzle(id string) LivePuzzleBase

	AddRequest(string, Request) error

	AddCallback(string, func(*Update)) error
}

// PuzzlePool represents the pool of interactable puzzles
type PuzzlePool struct {
	puzzles    map[string]LivePuzzleBase
	numWorkers int
}

// GlobalPuzzlePool represents all the puzzles
var GlobalPuzzlePool PuzzlePoolBase

// InitPuzzlePool assigns value to globalUserPool
func InitPuzzlePool() {
	GlobalPuzzlePool = &PuzzlePool{}
}

// AddPuzzle adds a puzzle to the pool
func (p *PuzzlePool) AddPuzzle(l LivePuzzleBase) {
	p.puzzles[l.ID()] = l
}

// GetPuzzle gets a puzzle from the pool
func (p *PuzzlePool) GetPuzzle(id string) LivePuzzleBase {
	return p.puzzles[id]
}

// AddRequest adds a web request to be processed by the workers
func (p *PuzzlePool) AddRequest(puzzleID string, r Request) error {
	puzzle, exists := p.puzzles[puzzleID]
	if exists {
		puzzle.AddRequest(r)
	}
	return fmt.Errorf("puzzle not in pool")
}

// AddCallback adds a callback to a puzzle
func (p *PuzzlePool) AddCallback(puzzleID string, cb func(*Update)) error {
	puzzle, exists := p.puzzles[puzzleID]
	if exists {
		puzzle.AddCallback(cb)
	}
	return fmt.Errorf("puzzle not in pool")
}
