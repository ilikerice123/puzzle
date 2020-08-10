package puzzle

import (
	"sync"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Base is an interface for the base Puzzle object
type Base interface {
	Do(r Request) bool

	GetUpdateBatch(position int, batch int) []Update

	// GetUpdate blocking gets a single update
	GetUpdate(position int) Update

	Complete() bool

	LastUpdated() time.Time
}

// Puzzle representing a non-threadsafe puzzle objec that implements Base
type Puzzle struct {
	id            string
	pieces        [][]*Piece
	heldPieces    map[string]*Piece
	size          int
	piecesCorrect int
	latestMove    int
	updates       []Update
	lock          *sync.Mutex
	// used for broadcasting to updates
	cv *sync.Cond
}

// NewPuzzle creates the new puzzle
func NewPuzzle(file string, xSize int, ySize int) *Puzzle {
	return &Puzzle{}
}

// Do does the request on the puzzle
func (p *Puzzle) Do(r Request) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	defer p.cv.Broadcast()

	piece := p.pieces[r.PieceX][r.PieceY]

	if piece.Metadata.HeldBy != "" && piece.Metadata.HeldBy != r.UserID {
		// piece is being held by someone else, no-op
		return false
	}

	if p.heldPieces[r.UserID] == nil {
		// hold, update held pieces
		piece.Metadata.HeldBy = r.UserID
		p.heldPieces[r.UserID] = piece
		p.updates = append(
			p.updates,
			Update{Action: HOLD, UserID: r.UserID, Piece1ID: piece.ID})
		return true
	}

	// swap (or release if same as held piece)
	otherPiece := p.heldPieces[r.UserID]
	otherPiece.Metadata.HeldBy = ""
	piece.Metadata.HeldBy = ""
	delta := p.swap(piece, otherPiece)
	p.piecesCorrect += delta

	p.updates = append(
		p.updates,
		Update{
			Action:   SWAP,
			UserID:   r.UserID,
			Piece1ID: piece.ID,
			Piece2ID: otherPiece.ID,
			Delta:    delta,
		})
	return true
}

// GetUpdateBatch blocking gets a batch of updates, up to num updates
func (p *Puzzle) GetUpdateBatch(position int, batch int) []Update {
	for position >= len(p.updates) {
		p.cv.L.Lock()
		p.cv.Wait()
		p.cv.L.Unlock()
	}
	cutoff := min(len(p.updates), position+batch)
	return p.updates[position:cutoff]
}

// GetUpdate blocking gets a single update
func (p *Puzzle) GetUpdate(position int) Update {
	for position >= len(p.updates) {
		p.cv.L.Lock()
		p.cv.Wait()
		p.cv.L.Unlock()
	}
	return p.updates[position]
}

// Complete returns if puzzle is finished
func (p *Puzzle) Complete() bool {
	return p.size == p.piecesCorrect
}

// LastUpdated returns when the puzzle was last updated
func (p *Puzzle) LastUpdated() time.Time {
	return time.Now()
}

// swap swaps piece1 and 2, and returns change in how many pieces are correct
func (p *Puzzle) swap(piece1 *Piece, piece2 *Piece) int {
	delta := 0
	if piece1.Correct() {
		delta--
	}
	if piece2.Correct() {
		delta--
	}
	piece2Pos := piece2.CurrPos
	// changing puzzle state
	p.pieces[piece2Pos.Y][piece2Pos.X] = piece1
	p.pieces[piece1.CurrPos.Y][piece1.CurrPos.X] = piece2
	// changing piece state
	piece2.CurrPos = piece1.CurrPos
	piece1.CurrPos = piece2Pos
	if piece1.CurrPos.Equals(piece1.DestPos) {
		delta++
	}
	if piece2.CurrPos.Equals(piece2.DestPos) {
		delta++
	}
	return delta
}
