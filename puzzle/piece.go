package puzzle

// Metadata contains all the information about a piece
type Metadata struct {
	ImgURL string
	HeldBy string
}

// Position represents a 2d position on the puzzle
type Position struct {
	X int
	Y int
}

// Piece represents a single puzzle piece
type Piece struct {
	DestPos  Position
	CurrPos  Position
	ID       int
	Metadata Metadata
}

// Equals compares different positions
func (p Position) Equals(other Position) bool {
	return p.X == other.X && p.Y == other.Y
}

// Correct returns if piece is in the correct position
func (p Piece) Correct() bool {
	return p.DestPos.Equals(p.CurrPos)
}
