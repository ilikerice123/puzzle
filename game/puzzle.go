package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ilikerice123/puzzle/picture"
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

	GetUpdate(position int) Update

	Shuffle()

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
	updates       []Update
	lock          sync.Locker
	cv            *sync.Cond
	users         UserPoolBase
}

// NewPuzzle creates the new puzzle from the file string of an image
func NewPuzzle(file string, ySize int, xSize int) (*Puzzle, error) {
	pieceNames, err := picture.SliceImage(file, ySize, xSize)
	if err != nil {
		return nil, err
	}
	puzzle := Puzzle{
		id:            uuid.New().String(),
		pieces:        make([][]*Piece, ySize),
		heldPieces:    make(map[string]*Piece),
		size:          ySize * xSize,
		piecesCorrect: 0,
		updates:       make([]Update, 0),
		lock:          &sync.Mutex{},
		cv:            sync.NewCond(&sync.Mutex{})}

	for i := range puzzle.pieces {
		puzzle.pieces[i] = make([]*Piece, xSize)
		for j := range puzzle.pieces[i] {
			puzzle.pieces[i][j] = &Piece{
				DestPos:  Position{Y: i, X: j},
				CurrPos:  Position{Y: i, X: j},
				ID:       i*xSize + j,
				Metadata: Metadata{ImgURL: pieceNames[i][j]}}
		}
	}
	puzzle.Shuffle()

	return &puzzle, nil
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

	reqUser := p.users.GetUser(r.UserID)
	if reqUser == nil {
		// somehow, it is not in pool
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

	// update user stats
	reqUser.PieceCount[p.id] = reqUser.PieceCount[p.id] + delta
	reqUser.LifetimePieces += delta

	p.updates = append(
		p.updates,
		Update{
			Action:   SWAP,
			UserID:   r.UserID,
			Piece1ID: piece.ID,
			Piece2ID: otherPiece.ID,
			Delta:    delta})
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

// Shuffle shuffles the puzzle
func (p *Puzzle) Shuffle() {

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
