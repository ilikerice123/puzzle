import React from 'react';
import { PuzzlePiece } from './game';
import * as CSS from 'csstype';

type PieceProps = {
    piece: PuzzlePiece
    host: string
    styles: CSS.Properties
}

export default function Piece(props: PieceProps) {
    return (
        <div>
            <img style={props.styles} src={`${props.host}/${props.piece.image}`} alt="" />
        </div>
    )
}