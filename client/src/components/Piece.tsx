import React from 'react';
import { PuzzlePieceObject, Pos } from './game';
import * as CSS from 'csstype';

type PieceProps = {
    piece: PuzzlePieceObject
    host: string
    styles: CSS.Properties
    onClick: (pos: Pos) => any
}

export default function PieceComponent(props: PieceProps) {
    return (
        <div>
            <img 
                draggable={false}
                style={props.styles} 
                src={`${props.host}/${props.piece.image}`} 
                alt=""
                onClick={() => props.onClick(props.piece.currPos)}
                />
        </div>
    )
}