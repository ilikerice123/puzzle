import React from 'react';
import PuzzleClient from './client'
import { RouteComponentProps, withRouter } from "react-router-dom";
import { PuzzleObject, PuzzlePieceObject } from './game';
import * as CSS from 'csstype';
import { debounce } from 'ts-debounce';
import Piece from './Piece';

// also servers as the interface for 'this.props.match.params'
interface PuzzleGameProps {
    puzzle: PuzzleObject
}

type PuzzleGameState = {
    puzzle: PuzzleObject
    maxPieceHeight: number
    maxPieceWidth: number
}

interface PieceLimits{
    maxHeight: number
    maxWidth: number
}

export default class PuzzleGameComponent extends React.Component<PuzzleGameProps, PuzzleGameState> {
    client = new PuzzleClient()
    resizeCallback: any

    constructor(props: PuzzleGameProps) {
        super(props);
        let {maxHeight, maxWidth} = this.getPieceLimits()
        this.state = {puzzle: props.puzzle, maxPieceHeight: maxHeight, maxPieceWidth: maxWidth}
    }
    
    componentDidMount() {
        this.resizeCallback = debounce(() => this.windowSizeChanged(), 1000)
        window.addEventListener("resize", this.resizeCallback)
    }
    
    componentWillUnmount() {
        window.removeEventListener('resize', this.resizeCallback);
    }
    
    windowSizeChanged() {
        let {maxHeight, maxWidth} = this.getPieceLimits()
        this.setState((prevState : PuzzleGameState) => {
            return {puzzle: prevState.puzzle, maxPieceHeight: maxHeight, maxPieceWidth: maxWidth}
        })
    }
    
    getPieceLimits(): PieceLimits {
        let width = window.innerWidth
        let height = window.innerHeight
        return {maxHeight: height/this.props.puzzle.ySize, maxWidth: width/this.props.puzzle.xSize}
    }

    render() {
        var pieceStyle : CSS.Properties = {
            maxWidth: `${this.state.maxPieceWidth.toString()}px`,
            maxHeight: `${this.state.maxPieceHeight.toString()}px`,
            border: "1px solid red"
        }
        
        return (
            <div>
                {this.state.puzzle != null && (
                    <table>
                        {this.state.puzzle.pieces.map((row: PuzzlePieceObject[]) => {
                            return (
                                <tr>
                                    {row.map((piece: PuzzlePieceObject) => {
                                        return (
                                            <td>
                                                <Piece 
                                                    piece={piece} 
                                                    host={this.client.host()}
                                                    styles={pieceStyle}
                                                /> 
                                            </td>
                                        )
                                    })}
                                </tr>
                            )
                        })}
                    </table>
                )}
            </div>
        )
    }
}

