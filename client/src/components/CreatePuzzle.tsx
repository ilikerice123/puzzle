import React from 'react';
import PuzzleClient from './client'
import { RouteComponentProps, withRouter } from 'react-router-dom';

interface CreatePuzzleProps extends RouteComponentProps<any> {
}

type CreatePuzzleState = {
    imageID?: string,
    ySize?: number,
    xSize?: number,
    loading: boolean
}

type ImageResponse = {
    uuid: string
}

class CreatePuzzle extends React.Component<CreatePuzzleProps, CreatePuzzleState> {
    client = new PuzzleClient()
    history: any

    constructor(props: CreatePuzzleProps){
        super(props)
        this.state = {loading: false}
    }

    async uploadFile(files: FileList | null) {
        if (files == null || files.length == 0) {
            console.log("no files selected!!")
            return
        }
        this.setState({loading: true})
        let res = await this.client.postFile<ImageResponse>("/images", "image", files[0])
        if (res == null || res.status != 200) {
            this.setState({loading: false, imageID: undefined})
            return
        }
        this.setState({loading: false, imageID: res.data.uuid})
    }

    async createPuzzle() {
        let id = this.state.imageID
        let ySize = this.state.ySize
        let xSize = this.state.xSize
        if (id == null || ySize == null || xSize == null) {
            return
        }
        let res = await this.client.postJson(`/puzzles/${id}`, {ySize: ySize, xSize: xSize})
        if (res == null || res.status != 200) {
            console.log("error occurred!")
            return
        }
        this.props.history.push(`/puzzles/${id}`)
    }

    render() {
        return (
            <div>
                Upload Puzzle: 
                <input 
                    type="file" 
                    accept="image/png, image/jpeg, image/gif"
                    onChange={(e) => this.uploadFile(e.target.files)} />
                X:
                <input
                    type="number" 
                    onChange={(e) => this.setState({xSize: parseInt(e.target.value)})} />
                Y:
                <input
                    type="number" 
                    onChange={(e) => this.setState({ySize: parseInt(e.target.value)})} />
                <button
                    type="button"
                    onClick={() => this.createPuzzle()}>
                    GO
                </button>
                {this.state.loading && <p>loading...</p>}
                <br />
                {this.state.imageID != null && (
                    <div>
                        <p>Preview:</p>
                        <img src={`${this.client.host()}/images/${this.state.imageID}/preview.jpeg`} alt="" />
                    </div>
                )}
                
            </div>
        )
    }
}

export default withRouter(CreatePuzzle)