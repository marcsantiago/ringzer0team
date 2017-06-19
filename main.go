package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var sudoku = `
+---+---+---+---+---+---+---+---+---+
| 6 |   |   |   | 1 | 2 |   | 3 |   |
+---+---+---+---+---+---+---+---+---+
|   | 1 | 2 |   | 3 |   |   | 4 | 7 |
+---+---+---+---+---+---+---+---+---+
| 9 |   | 5 |   | 4 | 7 |   | 1 |   |
+---+---+---+---+---+---+---+---+---+
|   |   |   | 1 | 2 | 9 | 3 | 5 |   |
+---+---+---+---+---+---+---+---+---+
| 1 | 2 | 9 | 3 |   |   |   |   | 8 |
+---+---+---+---+---+---+---+---+---+
|   |   | 6 | 4 | 7 | 8 |   | 2 |   |
+---+---+---+---+---+---+---+---+---+
|   |   |   |   |   |   | 5 |   |   |
+---+---+---+---+---+---+---+---+---+
|   | 9 | 3 |   |   | 4 |   | 8 | 1 |
+---+---+---+---+---+---+---+---+---+
|   | 6 | 4 |   | 8 | 1 | 2 |   | 3 |
+---+---+---+---+---+---+---+---+---+
`

const lBreak = "+---+---+---+---+---+---+---+---+---+"

func loadBoard(rawBoard io.Reader) [][]string {
	var buf bytes.Buffer
	io.Copy(&buf, rawBoard)
	var board [][]string
	for _, line := range strings.Split(buf.String(), lBreak) {
		if strings.TrimSpace(line) != "" {
			parts := strings.Split(line, "|")
			trimedEdges := parts[1 : len(parts)-1]
			var tmp []string
			for _, num := range trimedEdges {
				if strings.TrimSpace(num) == "" {
					tmp = append(tmp, "-1")

				} else {
					tmp = append(tmp, num)
				}
			}
			board = append(board, tmp)
		}
	}
	return board
}

func main() {

	for _, b := range loadBoard(strings.NewReader(sudoku)) {
		fmt.Printf("%+v length %d\n", b, len(b))
	}
}
