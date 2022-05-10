package main

import (
	"fmt"
	"snake/db"
	"snake/game"
	"strings"
)

var Initial = []byte{
	157, 79, 170, 8, 108, 234, 163, 16, 251, 181, 23, 148, 55, 162, 211, 186, 194, 222, 152, 207, 57, 97, 87, 45, 245, 141, 142, 40, 13, 92, 89, 64, 191, 102, 247, 178, 28, 138, 118, 68, 226, 24, 151, 103, 15, 139, 154, 244, 180, 83, 82, 196, 171, 167, 31, 155, 63, 246, 38, 200, 228, 120, 218, 204, 10, 238, 47, 56, 146, 185, 172, 158, 133, 53, 117, 42, 193, 241, 206, 86, 161, 0, 77, 243, 149, 239, 121, 129, 2, 85, 159, 59, 96, 164, 81, 220, 114, 18, 214, 65, 60, 125, 188, 201, 104, 174, 153, 75, 240, 223, 126, 35, 189, 113, 27, 236, 122, 143, 124, 73, 227, 43, 49, 67, 187, 48, 99, 250, 39, 20, 165, 115, 1, 177, 93, 232, 202, 249, 116, 54, 6, 242, 252, 69, 255, 22, 176, 197, 110, 5, 61, 169, 254, 183, 19, 229, 109, 150, 111, 131, 156, 253, 208, 145, 58, 179, 76, 7, 91, 78, 37, 233, 212, 9, 215, 192, 62, 209, 33, 32, 198, 168, 17, 195, 136, 166, 98, 130, 71, 248, 90, 217, 25, 30, 112, 34, 231, 3, 237, 21, 80, 224, 100, 66, 52, 84, 106, 4, 101, 205, 26, 105, 128, 225, 210, 135, 137, 175, 95, 70, 132, 203, 182, 29, 219, 190, 199, 44, 235, 140, 147, 74, 144, 46, 123, 216, 221, 14, 94, 127, 119, 36, 184, 88, 107, 12, 41, 134, 213, 72, 173, 160, 50, 51, 11, 230,
}
var Secret = []byte("ababubazzz")
var Counter = uint64(123)

func drawLevel(level *game.Level) {
	var lines [][]string

	for y := int64(0); y < level.Field.Height(); y++ {
		var line []string

		for x := int64(0); x < level.Field.Width(); x++ {
			cell, _ := level.Field.Get(x, y)

			switch cell {
			case game.CELL_EMPTY:
				line = append(line, ".")
				break
			case game.CELL_FOOD:
				line = append(line, "*")
				break
			}
		}

		lines = append(lines, line)
	}

	for _, coordinates := range level.Snake.Body() {
		if level.Field.Has(coordinates.X, coordinates.Y) {
			lines[coordinates.Y][coordinates.X] = "#"
		}
	}

	head, _ := level.Snake.Head()

	if level.Field.Has(head.X, head.Y) {
		lines[head.Y][head.X] = "@"
	}

	for _, line := range lines {
		fmt.Println(strings.Join(line, " "))
	}
}

func main() {
	db.Migrate()
	seed, err := GenerateSeed(Initial, Secret, Counter)
	if err != nil {
		fmt.Println(err)
		return
	}

	level, err := GenerateLevel(seed)
	if err != nil {
		fmt.Println(err)
		return
	}

	var direction game.Direction

	for i := 0; ; i++ {
		err := level.Step(direction)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(i, level.Status())
		drawLevel(level)

		if level.Status() != game.STATUS_UNFINISHED {
			break
		}

		var line string
		fmt.Scanln(&line)

		switch line {
		case "a":
			direction = game.DIRECTION_LEFT
			break

		case "d":
			direction = game.DIRECTION_RIGHT
			break

		case "w":
			direction = game.DIRECTION_UP
			break

		case "s":
			direction = game.DIRECTION_DOWN
			break
		}
	}
}
