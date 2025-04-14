package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var size int
	const defaultValue = 8

	fmt.Printf("Enter chessboard size (default %d): ", defaultValue)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		size = defaultValue
	} else {
		size, _ = strconv.Atoi(input)
	}

	fmt.Print(generateChessboard(size))
}

func generateChessboard(size int) string {
	var chessBoard strings.Builder

	const black string = "#"
	const white string = "_"
	const separator string = "|"

	for x := 1; x <= size; x++ {
		chessBoard.WriteString(separator)

		for y := 1; y <= size; y++ {
			switch (x + y) % 2 {
			case 0:
				chessBoard.WriteString(black)
			case 1:
				chessBoard.WriteString(white)
			}

			chessBoard.WriteString(separator)
		}

		chessBoard.WriteString("\n")
	}

	chessBoard.WriteString(separator)

	return chessBoard.String()
}
