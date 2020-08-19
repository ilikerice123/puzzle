
export interface Pos {
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
    imageWidth: number
    imageHeight: number
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

export interface PuzzleUpdateObject {
    id: number
    action: number
    userID: string
    piece1Pos: Pos
    piece2Pos: Pos
    delta: number
}

export interface PuzzleRequestObject {
    action: number
    userID: string
    position: Pos
}