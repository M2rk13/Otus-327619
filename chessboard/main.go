package main

import "fmt"

func main() {
	var size = 8

	chessBoard := generateChessboard(size)
	fmt.Print(chessBoard)
}

func generateChessboard(size int) string {
	var chess string
	var cell string

	const black string = "#"
	const white string = "_"
	const separator string = "|"

	for x := 1; x <= size; x++ {
		chess = chess + separator

		for y := 1; y <= size; y++ {
			switch (x + y) % 2 {
			case 0:
				cell = black
			case 1:
				cell = white
			}

			chess = chess + cell + separator
		}

		chess = chess + "\n"
	}

	return chess + separator
}
