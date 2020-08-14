import React from 'react';
import PuzzleClient from './client'
import { RouteComponentProps, withRouter } from "react-router-dom"
import PuzzleGameComponent from './PuzzleGame' 
import { timingSafeEqual } from 'crypto';
import { PuzzleObject } from './game';

// also servers as the interface for 'this.props.match.params'
type PuzzleProps = {
    id: string,
}

type PuzzleState = {
    puzzle?: PuzzleObject
}

class PuzzleImpl extends React.Component<PuzzleProps, PuzzleState> {
    client = new PuzzleClient()

    constructor(props: PuzzleProps){
        super(props)
        this.state = {puzzle: undefined}
    }

    async componentDidMount() {
        let res = await this.client.get<PuzzleObject>(`/puzzles/${this.props.id}`)
        if (res == null || res.status != 200) {
            console.log("error occurred!")
            return
        }
        this.setState({puzzle: res.data})
    }

    render() {
        return (
            <div>
                <p>Puzzle:</p>
                <pre>{JSON.stringify(this.state.puzzle, null, 4)}</pre>
                {this.state.puzzle != null && <PuzzleGameComponent puzzle={this.state.puzzle}/>}
            </div>
        )
    }
}


// used to inject id prop
export default withRouter((props: RouteComponentProps<PuzzleProps>) => (
    <PuzzleImpl id={props.match.params.id}/>
))