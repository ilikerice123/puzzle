import axios, { AxiosResponse } from 'axios'

export default class PuzzleClient {
    serverHost: string

    constructor(){
        this.serverHost = GetHost()
    }

    host(): string {
        return this.serverHost
    }

    // awaitable that doesn't throw
    async postFile<T>(url: string, name: string, file: File): Promise<AxiosResponse<T> | null> {
        var form = new FormData();
        form.append(name, file)
        let headers = {"Content-Type": "multipart/form-data"}
        try {
            return await axios.post(`${this.serverHost}${url}`, form, {headers: headers})
        } catch (e) {
            console.log("error occurred")
            console.log(e)
            return null
        }
    }
    
    // awaitable that doesn't throw
    async postJson<T>(url: string, data: object): Promise<AxiosResponse<T> | null> {
        let headers = {"Content-Type": "application/json"}
        try {
            return await axios.post(`${this.serverHost}${url}`, data, {headers: headers})
        } catch (e) {
            console.log("error occurred")
            console.log(e)
            return null
        }
    }

    async get<T>(url: string): Promise<AxiosResponse<T> | null> {
        try {
            return await axios.get<T>(`${this.serverHost}${url}`)
        } catch (e) {
            console.log("error occurred")
            console.log(e)
            return null
        }
    }
}

export function GetHost(): string {
    let url = new URL(window.location.href)
    if (url.host.includes("localhost")) {
        return "http://localhost:8000/api"
    } else {
        throw new Error("unimplemented!")
    }
}