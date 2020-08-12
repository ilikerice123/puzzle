package game

import (
	"fmt"
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

// PuzzleBase is an interface for the base Puzzle object
type PuzzleBase interface {
	Do(r Request) (bool, error)

	Shuffle()

	Complete() bool

	LastUpdatedTime() time.Time

	GetID() string
}

// Puzzle representing a non-threadsafe puzzle objec that implements PuzzleBase
type Puzzle struct {
	ID            string            `json:"id"`
	Pieces        [][]*Piece        `json:"pieces"`
	HeldPieces    map[string]*Piece `json:"held_pieces"`
	Size          int               `json:"size"`
	PiecesCorrect int               `json:"pieces_correct"`
	NextUpdateID  int               `json:"next_update_id"`
	XSize         int               `json:"x_size"`
	YSize         int               `json:"y_size"`
	LastUpdated   time.Time         `json:"last_updated"`
	updates       chan<- *Update
	users         UserPoolBase
}

// NewPuzzle creates the new puzzle from the file string of an image
func NewPuzzle(
	file string,
	ySize int,
	xSize int,
	updatesChannel chan<- *Update,
	users UserPoolBase) *Puzzle {
	pieceNames, err := picture.SliceImage(file, ySize, xSize)
	if err != nil {
		return nil
	}

	puzzle := Puzzle{
		ID:            uuid.New().String(),
		Pieces:        make([][]*Piece, ySize),
		HeldPieces:    make(map[string]*Piece),
		Size:          ySize * xSize,
		PiecesCorrect: 0,
		NextUpdateID:  0,
		LastUpdated:   time.Now(),
		updates:       updatesChannel,
		users:         users,
		XSize:         xSize,
		YSize:         ySize}

	for i := range puzzle.Pieces {
		puzzle.Pieces[i] = make([]*Piece, xSize)
		for j := range puzzle.Pieces[i] {
			puzzle.Pieces[i][j] = &Piece{
				DestPos:  Position{Y: i, X: j},
				CurrPos:  Position{Y: i, X: j},
				ID:       i*xSize + j,
				Metadata: Metadata{ImgURL: pieceNames[i][j]}}
		}
	}
	puzzle.Shuffle()

	return &puzzle
}

// GetID returns id of puzzle for the interface
func (p *Puzzle) GetID() string {
	return p.ID
}

// Do does the request on the puzzle
func (p *Puzzle) Do(r Request) (bool, error) {
	if r.PieceY >= p.YSize || r.PieceX >= p.XSize {
		return false, fmt.Errorf("piece x and y out of bounds")
	}

	piece := p.Pieces[r.PieceX][r.PieceY]
	reqUser := p.users.GetUser(r.UserID)
	if reqUser == nil {
		return false, fmt.Errorf("user not registered in pool")
	}

	if piece.Metadata.HeldBy != "" && piece.Metadata.HeldBy != r.UserID {
		// piece is being held by someone else, no-op
		return false, nil
	}

	p.LastUpdated = time.Now()
	if p.HeldPieces[r.UserID] == nil {
		// hold, update held pieces
		piece.Metadata.HeldBy = r.UserID
		p.HeldPieces[r.UserID] = piece
		p.updates <- p.newUpdate(HOLD, r.UserID, piece.ID, -1, 0)
		return true, nil
	}

	// swap (or release if same as held piece)
	otherPiece := p.HeldPieces[r.UserID]
	delta := p.swap(piece, otherPiece)
	p.PiecesCorrect += delta

	// update user stats
	reqUser.PieceCount[p.ID] = reqUser.PieceCount[p.ID] + delta
	reqUser.LifetimePieces += delta

	p.updates <- p.newUpdate(SWAP, r.UserID, piece.ID, otherPiece.ID, delta)
	return true, nil
}

// Shuffle shuffles the puzzle
func (p *Puzzle) Shuffle() {

}

// Complete returns if puzzle is finished
func (p *Puzzle) Complete() bool {
	return p.Size == p.PiecesCorrect
}

// LastUpdatedTime returns when the puzzle was last updated
func (p *Puzzle) LastUpdatedTime() time.Time {
	return p.LastUpdated
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
	p.Pieces[piece2Pos.Y][piece2Pos.X] = piece1
	p.Pieces[piece1.CurrPos.Y][piece1.CurrPos.X] = piece2
	// changing piece state
	piece2.CurrPos = piece1.CurrPos
	piece1.CurrPos = piece2Pos
	// after they swap, they are no longer being held
	piece1.Metadata.HeldBy = ""
	piece2.Metadata.HeldBy = ""
	if piece1.CurrPos.Equals(piece1.DestPos) {
		delta++
	}
	if piece2.CurrPos.Equals(piece2.DestPos) {
		delta++
	}
	return delta
}

// newUpdate creates an update for the puzzle
func (p *Puzzle) newUpdate(
	action action,
	userID string,
	piece1 int,
	piece2 int,
	delta int) *Update {
	// follow the update id dictated by the puzzle
	p.NextUpdateID++
	return &Update{
		ID:       p.NextUpdateID - 1,
		Action:   action,
		UserID:   userID,
		Piece1ID: piece1,
		Piece2ID: piece2,
		Delta:    delta}
}
