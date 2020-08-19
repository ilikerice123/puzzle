import React from 'react';
import PuzzleClient from './client'
import { RouteComponentProps, withRouter } from "react-router-dom"
import { PuzzleObject, UserObject, PuzzleUpdateObject, PuzzlePieceObject, Pos, PuzzleRequestObject } from './game';
import { debounce } from 'ts-debounce';
import * as CSS from 'csstype';
import { CSS_COLORS } from './colors'
import PieceComponent from './Piece';

interface PuzzleProps {
    user: UserObject | null,
    id: string
}

interface RoutingPuzzleProps extends RouteComponentProps<PuzzleParams>{
    user: UserObject | null,
}

// also servers as the interface for 'this.props.match.params'
interface PuzzleParams {
    id: string,
}

type PuzzleState = {
    puzzle: PuzzleObject | null
    user: UserObject | null
    // Map<string, number>
    stats: any | null 
    pieceLimits: PieceLimits
    done: boolean
}

// response from server
type PuzzleResponse = {
    Puzzle: PuzzleObject
}

interface PieceLimits{
    maxHeight: number
    maxWidth: number
}

const SWAP: number = 0
const HOLD: number = 1
const JOIN: number = 2
const LEAVE: number = 3

class PuzzleImpl extends React.Component<PuzzleProps, PuzzleState> {
    client = new PuzzleClient()
    conn?: WebSocket
    queue: PuzzleUpdateObject[]
    resizeCallback: any

    constructor(props: PuzzleProps) {
        super(props)
        this.state = {
            puzzle: null, 
            user: null, 
            stats: null, 
            pieceLimits: this.getPieceLimits(),
            done: false
        }
        this.userChanged = this.userChanged.bind(this)
        this.onPieceClicked = this.onPieceClicked.bind(this)
        this.loadPuzzle = this.loadPuzzle.bind(this)
        this.updatePuzzle = this.updatePuzzle.bind(this)
        this.windowSizeChanged = this.windowSizeChanged.bind(this)
        this.getPieceLimits = this.getPieceLimits.bind(this)
        this.queue = []
    }

    componentDidUpdate(prevProps: PuzzleProps) {
        if (this.props.user?.id == null) {
            return
        }
        if (prevProps.user == null || prevProps.user.id != this.props.user.id) {
            this.userChanged(this.props.user)
        }
    }

    componentDidMount() {
        this.resizeCallback = debounce(() => this.windowSizeChanged(), 1000)
        window.addEventListener("resize", this.resizeCallback)
        // load the user immediately if we have user
        if (this.props.user != null) {
            this.userChanged(this.props.user)
        }
    }
    
    componentWillUnmount() {
        window.removeEventListener('resize', this.resizeCallback);
    }
    
    windowSizeChanged() {
        if (this.state.puzzle == null) {
            return
        }
        let limits = this.getPieceLimits(this.state.puzzle)
        this.setState((prevState : PuzzleState) => {
            return {...prevState, pieceLimits: limits}
        })
    }

    // TODO: whole getting piece limits logic is a bit weird, need to change
    // we shouldn't have to call this when object does not exist, but we do
    getPieceLimits(puzzle?: PuzzleObject): PieceLimits {
        if (puzzle == null) {
            return {maxHeight: 0, maxWidth: 0}
        }
        let width = window.innerWidth
        let height = window.innerHeight
        return {
            maxHeight: Math.min(puzzle.imageHeight/puzzle.ySize, Math.floor(height/puzzle.ySize)), 
            maxWidth: Math.min(puzzle.imageHeight/puzzle.xSize, Math.floor(width/puzzle.xSize))
        }
    }

    onPieceClicked(piecePos: Pos) {
        if (this.state.user == null || this.conn == null || this.state.done) {
            return
        }
        let update: PuzzleRequestObject = {
            action: HOLD,
            userID: this.state.user.id,
            position: piecePos
        }
        this.conn.send(JSON.stringify(update))
    }

    // when the user is changed, this means that the user has logged in, either for the first time,
    // or as a different user, so we need to redo our websocket connection
    userChanged(user: UserObject) {
        if (this.conn != null) {
            this.conn.close(1000, "user logged off")
            // set puzzle to null, so we can store our updates
            this.setState((prevState: PuzzleState) => {
                return {...prevState, puzzle: null}
            })
        }
        this.conn = this.client.websocket(`/puzzles/${this.props.id}/ws?user=${user.id}`)
        this.conn.onmessage = this.updatePuzzle
        this.conn.onopen = this.loadPuzzle
        this.setState((prevState: PuzzleState) => {
            return {...prevState, user: user}
        })
    }

    updatePuzzle(ev: MessageEvent) {
        let update: PuzzleUpdateObject = JSON.parse(ev.data)
        if (this.state.puzzle == null) {
            this.queue.push(update)
            return
        }
        if (this.queue.length != 0) {
            for (let update of this.queue) {
                if (update.id >= this.state.puzzle.nextUpdateID) {
                    this.applyUpdate(update)
                }
            }
            this.queue = []
        }
        this.applyUpdate(update)
    }

    async applyUpdate(update: PuzzleUpdateObject) {
        if (this.state.puzzle == null) {
            return
        }
        switch(update.action) {
            case SWAP: {
                let pos = update.piece1Pos
                let pos1 = update.piece2Pos
                let piece = {...this.state.puzzle.pieces[pos.Y][pos.X]}
                let piece1 = {...this.state.puzzle.pieces[pos1.Y][pos1.X]}
                piece.heldBy = ""
                piece1.heldBy = ""
                // since we are swapping the pieces, we have to swap the cur pos
                let tempPos = piece1.currPos
                piece1.currPos = piece.currPos
                piece.currPos = tempPos
                this.updatePuzzlePieces(pos, piece1, pos1, piece, update.delta)
                this.setState((prevState) => {
                    let stats = {...prevState.stats}
                    if (stats.hasOwnProperty(update.userID)) {
                        stats[update.userID] += update.delta
                    } else {
                        stats[update.userID] = update.delta
                    }
                    return ({...prevState, stats: stats})
                })
                return
            }
            case HOLD: {
                let pos = update.piece1Pos
                let piece = {...this.state.puzzle.pieces[pos.Y][pos.X]}
                piece.heldBy = update.userID
                this.updatePuzzlePieces(pos, piece)
                return
            }
            case JOIN: {
                let res = await this.client.get<UserObject>(`/users/${update.userID}`)
                if (res == null || res.data == null) {
                    return
                } 
                let users:any = {...this.state.puzzle.currentUsers}
                users[update.userID] = res.data
                this.updatePuzzleUsers(users)
                return
            }
            case LEAVE: {
                let users:any = {...this.state.puzzle.currentUsers}
                delete users[update.userID]
                this.updatePuzzleUsers(users)
                return
            }
        }
    }

    updatePuzzlePieces(
        pos: Pos, piece: PuzzlePieceObject, 
        pos1?: Pos, piece1?: PuzzlePieceObject, delta?: number
    ) {
        if (this.state.puzzle == null) {
            return
        }

        this.setState((prevState: PuzzleState) => {
            if (prevState.puzzle == null) {
                return prevState
            }
            // shallow copy of the outer array
            let pieceRows = [...prevState.puzzle.pieces]
            // shallow copy of the row that the piece belongs
            let pieceRow = [...pieceRows[pos.Y]]
            // update the piece
            pieceRow[pos.X] = piece
            // update the reference on pieceRows
            pieceRows[pos.Y] = pieceRow

            if (pos1 != null && piece1 != null) {
                if (pos1.Y == pos.Y) {
                    // we don't have to shallow copy the row again
                    pieceRow[pos1.X] = piece1
                } else {
                    let pieceRow1 = [...pieceRows[pos1.Y]]
                    pieceRow1[pos1.X] = piece1
                    pieceRows[pos1.Y] = pieceRow1
                }
            }

            let done = false
            if (prevState.puzzle.piecesCorrect + (delta || 0) == prevState.puzzle.size) {
                done = true
            }

            return {
                ...prevState,
                puzzle: {
                    ...prevState.puzzle,
                    pieces: pieceRows,
                    piecesCorrect: prevState.puzzle.piecesCorrect + (delta || 0)
                },
                done: done
            }
        })
    }

    updatePuzzleUsers(users: Map<string, UserObject>) {
        this.setState((prevState: PuzzleState) => {
            if (prevState.puzzle == null) {
                return prevState
            }
            return {
                ...prevState,
                puzzle: {
                    ...prevState.puzzle,
                    currentUsers: users
                }
            }
        })
    }

    async loadPuzzle() {
        let res = await this.client.get<PuzzleResponse>(`/puzzles/${this.props.id}`)
        if (res == null || res.status != 200) {
            console.log("error occurred!")
            return
        }
        this.setState((prevState: PuzzleState) => {
            // should not happen ever - just here to make typescript happy
            if (res == null) {
                return prevState
            }
            let limits = this.getPieceLimits(res.data.Puzzle)
            return {...prevState, puzzle: res.data.Puzzle, stats: {}, pieceLimits: limits}
        })
    }

    render() {
        return (
            <div>
                {this.state.user == null && <h2 style={{padding: "20px"}}>Sign in to start playing</h2>}
                {this.state.done && <h2> PUZZLE COMPLETE! </h2>}
                {this.state.puzzle != null && (
                    <PuzzleGameComponent 
                        host={this.client.host()}
                        puzzle={this.state.puzzle} 
                        stats={this.state.stats}
                        pieceClicked={this.onPieceClicked}
                        pieceLimits={this.state.pieceLimits}
                    />
                )}
            </div>
        )
    }
}

const PuzzleStyles = {
    table: {
        borderSpacing: "0px",
    },
    tr: {
        padding: "0px",
        lineHeight: "0px"
    },
    td: {
        padding: "0px",
        margin: "0px",
        lineHeight: "0px"
    },
    img: {
        padding: "0px",
        margin: "0px",
        outlineOffset: "-5px"
    }
}

function PuzzleGameComponent(
    props: {
        host: string,
        puzzle: PuzzleObject,
        // Map<string, number>
        stats: any,
        pieceClicked: (pos: Pos) => any,
        pieceLimits: PieceLimits
    }
) {
    var pieceStyle : CSS.Properties = {
        maxWidth: `${props.pieceLimits.maxWidth.toString()}px`,
        maxHeight: `${props.pieceLimits.maxHeight.toString()}px`,
    }
    
    return (
        <div>
            <div>
                <ul>
                    {Array.from(Object.values(props.puzzle.currentUsers)).map(
                        (user: UserObject) => (
                            <li 
                                key={user.id} 
                                style={{color: userColor(user.id)}}
                            >
                                {user.name} ({props.stats[user.id] || 0})
                            </li>
                        )
                    )}
                </ul>
                <table style={PuzzleStyles.table}>
                    <tbody>
                        {props.puzzle.pieces.map((row: PuzzlePieceObject[], idx: number) => {
                            return (
                                <tr key={idx} style={PuzzleStyles.tr}>
                                    {row.map((piece: PuzzlePieceObject, idx: number) => {
                                        return (
                                            <td key={idx} style={PuzzleStyles.td}>
                                                <PieceComponent
                                                    piece={piece} 
                                                    host={props.host}
                                                    styles={{
                                                        ...pieceStyle,
                                                        ...PuzzleStyles.img,
                                                        outline: `5px solid ${userColor(piece.heldBy)}`
                                                        // MozBoxShadow: `inset 0px 0px 0px 5px ${userColor(piece.heldBy)}`,
                                                        // boxShadow: `inset 0px 0px 0px 5px ${userColor(piece.heldBy)}`,
                                                        // WebkitBoxShadow: `inset 0px 0px 0px 5px ${userColor(piece.heldBy)}`
                                                    }}
                                                    onClick={props.pieceClicked}
                                                /> 
                                            </td>
                                        )
                                    })}
                                </tr>
                            )
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

// hash the id into number, then use number to index css color array
export function userColor(id: string): string{
    if (id == "") {
        return "Transparent"
    }

    var hash = 0, i, chr;
    for (i = 0; i < id.length; i++) {
      chr   = id.charCodeAt(i);
      hash  = ((hash << 5) - hash) + chr;
      hash |= 0; // Convert to 32bit integer
    }
    
    /* JavaScript does bitwise operations (like XOR, above) on 32-bit signed
        * integers. Since we want the results to be always positive, convert the
        * signed int to an unsigned by doing an unsigned bitshift. */
    return CSS_COLORS[(hash >>> 0)%(CSS_COLORS.length)]
}



// used to inject id prop
export default withRouter((props: RoutingPuzzleProps) => (
    <PuzzleImpl id={props.match.params.id} user={props.user}/>
))