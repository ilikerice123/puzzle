
type Pos = {
    X: number
    Y: number
}

export interface PuzzleGame {
    id: string
    pieces: PuzzlePiece[][]
    heldPieces: Map<string, string>
    size: number
    piecesCorrect: number
    nextUpdateID: number
    xSize: number
    ySize: number
    lastUpdated: string
    currentUsers: Map<string, PuzzleUser>
}

export interface PuzzleUser {
    id: string
    name: string
    created: string
    lifetimePieces: number
}

export interface PuzzlePiece {
    currPos: Pos
    image: string
    heldBy: string
}