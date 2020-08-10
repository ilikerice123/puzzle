package puzzle

import "time"

// Base is an interface for the base Puzzle object
type Base interface {
	Do(r Request) bool

	GetUpdateBatch(position int, num int) []Update

	// GetUpdate blocking gets a single update
	GetUpdate(position int) Update

	Complete() bool

	LastUpdated() time.Time
}

// Puzzle representing a non-threadsafe puzzle objec that implements Base
type Puzzle struct {
	pieces           [][]Piece
	heldPieces       map[string]*Piece
	size             int
	numPiecesCorrect int
	latestMove       int
}

// Do does the request on the puzzle
func (p Puzzle) Do(r Request) bool {
	return false
}

// GetUpdateBatch blocking gets a batch of updates, up to num updates
func (p Puzzle) GetUpdateBatch(position int, num int) []Update {
	return nil
}

// GetUpdate blocking gets a single update
func (p Puzzle) GetUpdate(position int) Update {
	return Update{}
}

// Complete returns if puzzle is finished
func (p Puzzle) Complete() bool {
	return p.size == p.numPiecesCorrect
}

// LastUpdated returns when the puzzle was last updated
func (p Puzzle) lastUpdated() time.Time {
	return time.Now()
}
