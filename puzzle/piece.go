package puzzle

// Metadata contains all the information about a piece
type Metadata struct {
	imgURL string
	heldBy string
}

// Piece represents a single puzzle piece
type Piece struct {
	destX    int
	destY    int
	id       int
	metadata Metadata
}
