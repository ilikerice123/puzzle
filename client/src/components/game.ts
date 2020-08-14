
type Pos = {
    X: number
    Y: number
}

export interface PuzzleObject {
    id: string
    pieces: PuzzlePieceObject[][]
    heldPieces: Map<string, string>
    size: number
    piecesCorrect: number
    nextUpdateID: number
    xSize: number
    ySize: number
    lastUpdated: string
    currentUsers: Map<string, UserObject>
}

export interface UserObject {
    id: string
    name: string
    created: string
    lifetimePieces: number
}

export interface PuzzlePieceObject {
    currPos: Pos
    image: string
    heldBy: string
}