package main

import (
	"os"

	"github.com/ilikerice123/puzzle/picture"
)

func main() {
	picture.SliceImage(os.Args[1], 4, 4)
}
