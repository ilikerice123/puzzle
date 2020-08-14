package game

// Position represents a 2d position on the puzzle
type Position struct {
	X, Y int
}

// Piece represents a single puzzle piece
type Piece struct {
	DestPos   Position `json:"-"`
	CurrPos   Position `json:"currPos"`
	ID        int      `json:"-"`
	ImageFile string   `json:"image"`
	HeldBy    string   `json:"heldBy"`
}

// Equals compares different positions
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Correct returns if piece is in the correct position
func (p Piece) Correct() bool {
	return p.DestPos.Equals(p.CurrPos)
}
