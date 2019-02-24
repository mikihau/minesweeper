package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var neighbors = [8]struct {
	x, y int
}{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

type cell struct {
	num      int
	revealed bool
	flagged  bool
	isMine   bool
}

type board struct {
	width, height, numMines int
	grid                    [][]cell
	numFlagged, numRevealed int
	active                  bool
	started                 time.Time
	gameDurationSeconds     float64
}

func (b *board) withinBoarder(x, y int) bool {
	return x >= 0 && x < b.height && y >= 0 && y < b.width
}

func (b *board) isWin() bool {
	return b.numFlagged == b.numMines && b.numRevealed == b.width*b.height-b.numMines
}

func (b *board) secondsElapsed() float64 {
	if !b.active {
		return b.gameDurationSeconds
	}
	if b.started.IsZero() {
		return float64(0)
	}
	return time.Since(b.started).Seconds()
}

func makeBoard(width, height, numMines int) *board {
	board := board{}
	board.active = true
	board.width, board.height = width, height
	board.grid = make([][]cell, height)
	for i := 0; i < height; i++ {
		board.grid[i] = make([]cell, width)
	}
	for board.numMines < numMines {
		x, y := rand.Intn(height), rand.Intn(width)
		if !board.grid[x][y].isMine {
			board.grid[x][y].isMine = true
			board.numMines++
			for _, neighbor := range neighbors {
				neighborX, neighborY := x+neighbor.x, y+neighbor.y
				if board.withinBoarder(neighborX, neighborY) {
					board.grid[neighborX][neighborY].num++
				}
			}
		}
	}
	return &board
}

func (b *board) explore(x, y int) {
	cell := &b.grid[x][y]
	if cell.revealed || cell.flagged || cell.isMine {
		return
	}
	if cell.num >= 0 {
		cell.revealed = true
		b.numRevealed++
	}
	if cell.num == 0 {
		for _, neighbor := range neighbors {
			neighborX, neighborY := x+neighbor.x, y+neighbor.y
			if b.withinBoarder(neighborX, neighborY) {
				b.explore(neighborX, neighborY)
			}
		}
	}
}

func (b *board) updateOnReveal(x, y int) {
	if !b.withinBoarder(x, y) || !b.active {
		return
	}
	if b.started.IsZero() {
		b.started = time.Now()
	}
	cell := &b.grid[x][y]
	if cell.revealed || cell.flagged {
		return
	}
	if cell.isMine {
		cell.revealed = true
		b.numRevealed++
		fmt.Println("Game over :(")
		b.gameDurationSeconds = b.secondsElapsed()
		b.active = false
	}
	b.explore(x, y)
	if b.isWin() {
		fmt.Println("You win!")
		b.gameDurationSeconds = b.secondsElapsed()
		b.active = false
	}
}

func (b *board) updateOnFlag(x, y int) {
	if !b.withinBoarder(x, y) || b.grid[x][y].revealed || !b.active {
		return
	}
	if b.started.IsZero() {
		b.started = time.Now()
	}
	if !b.grid[x][y].flagged {
		b.grid[x][y].flagged = true
		b.numFlagged++
	} else {
		b.grid[x][y].flagged = false
		b.numFlagged--
	}
	if b.isWin() {
		fmt.Println("You win!")
		b.gameDurationSeconds = b.secondsElapsed()
		b.active = false
	}
}

func (b *board) updateOnNewGame(width, height, numMines int) {
	*b = *makeBoard(width, height, numMines)
}

func pad(s string, numDigits int, char string) string {
	result := s
	for len(result) < numDigits {
		result = char + result
	}
	return result
}

func (b board) String() string {
	/*  [legend]
	flagged mine: #
	unrevealed: -
	a mine revealed: X
	*/
	digitsHorizontal, digitsVertical := len(strconv.Itoa(b.width-1)), len(strconv.Itoa(b.height-1))
	result := fmt.Sprintf("Mines remaining: %v\nTime elapsed: %.0f seconds\n", b.numMines-b.numFlagged, b.secondsElapsed())
	// horizontal indices
	result += strings.Repeat(" ", digitsVertical+1)
	for c := range b.grid[0] {
		result += fmt.Sprintf(" %v", pad(strconv.Itoa(c), digitsHorizontal, "0"))
	}
	result += "\n" + strings.Repeat(" ", digitsVertical+1) + strings.Repeat("-", b.width*(digitsHorizontal+1)) + "\n"
	// vertical indices before each row
	for r, row := range b.grid {
		result += fmt.Sprintf("%v|", pad(strconv.Itoa(r), digitsVertical, "0"))
		for _, cell := range row {
			switch {
			case cell.flagged:
				result += pad("#", digitsHorizontal+1, " ")
			case !cell.revealed:
				result += pad("-", digitsHorizontal+1, " ")
			case cell.revealed && cell.isMine:
				result += pad("X", digitsHorizontal+1, " ")
			default:
				result += pad(strconv.Itoa(cell.num), digitsHorizontal+1, " ")
			}
		}
		result += "\n"
	}
	return result
}

// for debugging
func (b board) showUnderlyingBoard() string {
	result := ""
	for _, row := range b.grid {
		for _, cell := range row {
			if cell.isMine {
				result += " #"
			} else {
				result += fmt.Sprintf(" %v", cell.num)
			}
		}
		result += "\n"
	}
	return result
}

func parseIntArgs(command string, numArgsExpected int) ([]int, error) {
	fields := strings.Fields(command)[1:]
	if len(fields) != numArgsExpected {
		message := fmt.Sprintf("Wrong number of arguments for command '%v': expecting %v, got %v -- type 'h' for help", command, numArgsExpected, len(fields))
		return nil, errors.New(message)
	}
	result := make([]int, numArgsExpected)
	for i, field := range fields {
		if num, err := strconv.Atoi(field); err == nil {
			result[i] = num
		} else {
			message := fmt.Sprintf("Expecting integers following command '%v'-- type 'h' for help", command)
			return nil, errors.New(message)
		}
	}
	return result, nil
}

func eval() {
	builtinBoardDimensions := map[string][]int{
		"default":      []int{8, 8, 10},
		"beginner":     []int{8, 8, 10},
		"intermediate": []int{16, 16, 40},
		"expert":       []int{31, 16, 99},
	}
	board := makeBoard(builtinBoardDimensions["default"][0], builtinBoardDimensions["default"][1], builtinBoardDimensions["default"][2])
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		if command := scanner.Text(); command != "" {
			switch {
			case command == "h" || command == "help":
				fmt.Println("Available commands:\nh(help) -- print help\nn(new) [beginner|intermediate|expert]|[width height numMines] -- start a new game\nr(reveal) <row> <col> -- reveal a cell\nf(flag) <row> <col> -- flag a cell")
			case command == "e" || command == "exit":
				return
			case strings.HasPrefix(command, "n") || strings.HasPrefix(command, "new"):
				fields := strings.Fields(command)[1:]
				var dimensions []int
				if len(fields) == 0 {
					dimensions = builtinBoardDimensions["default"]
				} else if len(fields) == 1 {
					var ok bool
					dimensions, ok = builtinBoardDimensions[fields[0]]
					if !ok {
						fmt.Println("Expecting either 'beginner', 'intermediate' or 'expert' -- type 'h' for help")
						break
					}
				} else {
					var err error
					dimensions, err = parseIntArgs(command, 3)
					if err != nil {
						fmt.Println(err)
						break
					}
				}
				board.updateOnNewGame(dimensions[0], dimensions[1], dimensions[2])
				fmt.Println(board)
			case strings.HasPrefix(command, "r ") || strings.HasPrefix(command, "reveal "):
				if args, err := parseIntArgs(command, 2); err != nil {
					fmt.Println(err)
				} else {
					board.updateOnReveal(args[0], args[1])
				}
				fmt.Println(board)
				if !board.active {
					fmt.Println("This game has ended -- type 'n' for a new game.")
				}
			case strings.HasPrefix(command, "f ") || strings.HasPrefix(command, "flag "):
				if args, err := parseIntArgs(command, 2); err != nil {
					fmt.Println(err)
				} else {
					board.updateOnFlag(args[0], args[1])
				}
				fmt.Println(board)
				if !board.active {
					fmt.Println("This game has ended -- type 'n' for a new game.")
				}
			default:
				fmt.Println("Bad command -- type 'h' for help")
			}
		}
		fmt.Print("> ")
	}
	if scanner.Err() != nil {
		fmt.Printf("Encountered error while reading command: %v; ignoring this command.\n", scanner.Err())
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	eval()
}
