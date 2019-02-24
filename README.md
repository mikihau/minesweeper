# Minesweeper
Minesweeper game in the command line, written in go.  
<img src="./demo.gif" alt="demo gif">

## Installing

```bash
git clone https://github.com/mikihau/minesweeper.git
```

## Running
```bash
$ go run minesweeper.go
> h
Available commands:
h(help) -- print help
n(new) [beginner|intermediate|expert]|[width height numMines] -- start a new game
r(reveal) <row> <col> -- reveal a cell
f(flag) <row> <col> -- flag a cell
> 
```

Example commands:
- `n intermediate`: starts a new game at intermediate level
- `n 8 8 10`: starts a new game with board size 8 * 8 and 10 mines
- `n`: starts a new game with default difficulty level
- `r 3 5`: reveal the cell on row 3, column 5
- `f 3 5`: flag the cell on row 3, column 5 as a mine
