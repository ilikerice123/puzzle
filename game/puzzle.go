package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ilikerice123/puzzle/fs"
	"github.com/ilikerice123/puzzle/store"

	"github.com/ilikerice123/puzzle/picture"
)

// PuzzleBase is an interface for the base Puzzle object
type PuzzleBase interface {
	Do(r Request) error

	Shuffle()

	Complete() bool

	LastUpdatedTime() time.Time

	GetID() string

	Results() map[string]int

	OnComplete()
}

// Puzzle representing a non-threadsafe puzzle objec that implements PuzzleBase
type Puzzle struct {
	ID            string                 `json:"id"`
	Pieces        [][]*Piece             `json:"pieces"`
	HeldPieces    map[string]*Piece      `json:"heldPieces"`
	Size          int                    `json:"size"`
	PiecesCorrect int                    `json:"piecesCorrect"`
	NextUpdateID  int                    `json:"nextUpdateID"`
	XSize         int                    `json:"xSize"`
	YSize         int                    `json:"ySize"`
	ImageWidth    int                    `json:"imageWidth"`
	ImageHeight   int                    `json:"imageHeight"`
	LastUpdated   time.Time              `json:"lastUpdated"`
	CurrentUsers  map[string]*store.User `json:"currentUsers"`
	updates       chan<- *Update
	users         UserPoolBase
}

// NewPuzzle creates the new puzzle from the file string of an image
func NewPuzzle(
	id string,
	file string,
	ySize int,
	xSize int,
	updatesChannel chan<- *Update,
	users UserPoolBase) *Puzzle {
	pieceNames, err := picture.SliceImage(file, ySize, xSize)
	if err != nil {
		return nil
	}

	// duplicate call to normalizeImage and loadImage just to get the dimensions
	image, _ := fs.LoadImage(file)
	imageWidth, imageHeight := picture.NormalizeImage(image, ySize, xSize)

	puzzle := Puzzle{
		ID:            id,
		Pieces:        make([][]*Piece, ySize),
		HeldPieces:    make(map[string]*Piece),
		Size:          ySize * xSize,
		PiecesCorrect: 0,
		NextUpdateID:  0,
		LastUpdated:   time.Now(),
		CurrentUsers:  make(map[string]*store.User),
		updates:       updatesChannel,
		users:         users,
		XSize:         xSize,
		YSize:         ySize,
		ImageWidth:    imageWidth,
		ImageHeight:   imageHeight,
	}

	for i := range puzzle.Pieces {
		puzzle.Pieces[i] = make([]*Piece, xSize)
		for j := range puzzle.Pieces[i] {
			puzzle.Pieces[i][j] = &Piece{
				DestPos:   Position{Y: i, X: j},
				CurrPos:   Position{Y: i, X: j},
				ID:        i*xSize + j,
				HeldBy:    "",
				ImageFile: pieceNames[i][j]}
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
func (p *Puzzle) Do(r Request) error {
	if p.PiecesCorrect == p.Size {
		return fmt.Errorf("puzzle complete")
	}

	switch r.Action {
	case HOLD:
		return p.hold(r)
	case JOIN:
		return p.addUser(r.UserID)
	case LEAVE:
		p.removeUser(r.UserID)
		return nil
	default:
		return fmt.Errorf("unknown action")
	}
}

// Shuffle shuffles the puzzle based on Fisher–Yates shuffle, modified so
// piecesCorrect == 0 every time
func (p *Puzzle) Shuffle() {
	for i := range p.Pieces {
		for j := range p.Pieces[i] {
			if i == len(p.Pieces)-1 && j == len(p.Pieces[i])-1 {
				return
			}
			currIdx := i*(p.XSize) + j
			idx := rand.Intn((p.Size)-currIdx-1) + currIdx + 1

			// log.Printf("swapping (%d,%d) with (%d,%d)\n", i, j, k, l)
			p.swap(p.Pieces[i][j], p.Pieces[idx/p.XSize][idx%p.XSize])
		}
	}
}

// OnComplete executes completion logic
func (p *Puzzle) OnComplete() {

}

// Results returns the results of the puzzle
func (p *Puzzle) Results() map[string]int {
	results := make(map[string]int)
	for userID, user := range p.CurrentUsers {
		results[userID] = user.PieceCount[p.ID]
	}
	return results
}

// Complete returns if puzzle is finished
func (p *Puzzle) Complete() bool {
	return p.Size == p.PiecesCorrect
}

// LastUpdatedTime returns when the puzzle was last updated
func (p *Puzzle) LastUpdatedTime() time.Time {
	return p.LastUpdated
}

func (p *Puzzle) hold(r Request) error {
	if r.PiecePos.Y >= p.YSize || r.PiecePos.X >= p.XSize {
		return fmt.Errorf("piece x and y out of bounds")
	}
	reqUser, exists := p.CurrentUsers[r.UserID]
	if !exists {
		return fmt.Errorf("puzzle's current users doesn't include user id")
	}

	piece := p.Pieces[r.PiecePos.Y][r.PiecePos.X]
	if piece.HeldBy != "" && piece.HeldBy != r.UserID {
		// piece is being held by someone else, no-op
		return nil
	}

	p.LastUpdated = time.Now()
	if p.HeldPieces[r.UserID] == nil {
		// hold, update held pieces
		piece.HeldBy = r.UserID
		p.HeldPieces[r.UserID] = piece
		p.updates <- p.newUpdate(HOLD, r.UserID, piece.CurrPos, Position{}, 0)
		return nil
	}

	otherPiece := p.HeldPieces[r.UserID]
	delta := p.swap(piece, otherPiece)
	// swap (or release if same as held piece)
	// TODO: needs to be after swap for some reason, or else the pointer is gone? what?
	delete(p.HeldPieces, r.UserID)
	p.PiecesCorrect += delta
	if p.Complete() {
		p.OnComplete()
	}

	// update user stats
	reqUser.PieceCount[p.ID] = reqUser.PieceCount[p.ID] + delta
	reqUser.LifetimePieces += delta

	p.updates <- p.newUpdate(SWAP, r.UserID, piece.CurrPos, otherPiece.CurrPos, delta)
	return nil
}

// addUser adds a user to current users
func (p *Puzzle) addUser(id string) error {
	u := p.users.GetUser(id)
	if u == nil {
		return fmt.Errorf("user not registered in pool")
	}

	if _, exists := p.CurrentUsers[u.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	p.CurrentUsers[u.ID] = u
	p.updates <- p.newUpdate(JOIN, id, Position{}, Position{}, 0)
	return nil
}

// removeUser removes a user from current users
func (p *Puzzle) removeUser(id string) {
	delete(p.CurrentUsers, id)
	if heldPiece, exists := p.HeldPieces[id]; exists {
		heldPiece.HeldBy = ""
		p.updates <- p.newUpdate(SWAP, id, heldPiece.CurrPos, heldPiece.CurrPos, 0)
	}
	delete(p.HeldPieces, id)
	p.updates <- p.newUpdate(LEAVE, id, Position{}, Position{}, 0)
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
	piece1.HeldBy = ""
	piece2.HeldBy = ""

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
	piece1 Position,
	piece2 Position,
	delta int) *Update {
	p.NextUpdateID++
	return &Update{
		ID:        p.NextUpdateID - 1,
		Action:    action,
		UserID:    userID,
		Piece1Pos: piece1,
		Piece2Pos: piece2,
		Delta:     delta}
}
