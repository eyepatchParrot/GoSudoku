package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	grid [][]int
}

func NewPuzzle() (p *Puzzle) {
	p = new(Puzzle)
	p.grid = make([][]int, 9)
	for i := 0; i < 9; i++ {
		p.grid[i] = make([]int, 9)
	}
	return p
}

func (p *Puzzle) Load(src string) (err error) {
	src = deleteWhitespace(src)
	if len(src) != 81 {
		return errors.New("wrong number of numbers/spaces in file")
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			i := (r * 9) + c
			var curNum int
			if src[i] == '_' {
				curNum = 0
			} else {
				curNum, err = strconv.Atoi(string(src[i]))
				assertOk(err)
			}
			p.grid[r][c] = curNum
		}
	}
	return nil
}

func (p *Puzzle) printBorder() {
	fmt.Println("+-----+-----+-----+")
}

func (p *Puzzle) printRow(row int) {
	printLine := ""
	for c := 0; c < 9; c++ {
		if (c % 3) == 0 {
			printLine = printLine + "|"
		}
		if p.grid[row][c] == 0 {
			printLine = printLine + "_"
		} else {
			printLine = printLine + strconv.Itoa(p.grid[row][c])
		}
		if (c % 3) != 2 {
			printLine = printLine + " "
		}
	}
	printLine = printLine + "|"
	fmt.Println(printLine)
}

func (p *Puzzle) Print() {
	for r := 0; r < 9; r++ {
		if (r % 3) == 0 {
			p.printBorder()
		}
		p.printRow(r)
	}
	p.printBorder()
}

// assumes that all numbers are valid
func (p *Puzzle) isComplete() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if p.grid[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

func (p *Puzzle) getNextSpaceAfter(startR, startC int) (r, c int, err error) {
	r, c = startR, startC
	for r = startR; r < 9; r++ {
		if r == startR {
			c = startC
		} else {
			c = 0
		}
		for ; c < 9; c++ {
			if p.grid[r][c] == 0 {
				return r, c, nil
			}
		}
	}
	return -1, -1, errors.New("getNextSpace called without an empty space")
}

func (p *Puzzle) boxIsOk(r, c, n int) bool {
	boxR, boxC := r-(r%3), c-(c%3)
	for testR := boxR; testR < boxR+3; testR++ {
		for testC := boxC; testC < boxC+3; testC++ {
			if p.grid[testR][testC] == n {
				return false
			}
		}
	}
	return true
}

func (p *Puzzle) rowIsOk(r, c, n int) bool {
	for testC := 0; testC < 9; testC++ {
		if p.grid[r][testC] == n {
			return false
		}
	}
	return true
}

func (p *Puzzle) colIsOk(r, c, n int) bool {
	for testR := 0; testR < 9; testR++ {
		if p.grid[testR][c] == n {
			return false
		}
	}
	return true
}

func (p *Puzzle) numberFits(r, c, n int) bool {
	if p.boxIsOk(r, c, n) && p.rowIsOk(r, c, n) && p.colIsOk(r, c, n) {
		return true
	} else {
		return false
	}
	return true
}

func (p *Puzzle) Copy() (ret *Puzzle) {
	ret = NewPuzzle()
	for r := 0; r < 9; r++ {
		copy(ret.grid[r], p.grid[r])
	}
	return ret
}

func (p *Puzzle) findBestSpace() (bestR, bestC, bestNum int) {
	bestR, bestC, err := p.getNextSpaceAfter(0, 0)
	if err != nil {
		panic(errors.New("findBestSpace called without a new space"))
	}
	bestNum = 10
	for r, c := bestR, bestC; err == nil; r, c, err = p.getNextSpaceAfter(r+1, c+1) {
		numHits := 0
		for i := 1; i < 10; i++ {
			if p.numberFits(r, c, i) {
				numHits++
			}
		}
		if numHits < bestNum {
			bestR, bestC, bestNum = r, c, numHits
			if numHits == 0 {
				break
			}
		}
	}
	return
}

func (p *Puzzle) Solve(undos *int) *Puzzle {
	if p.isComplete() {
		return p
	}

	r, c, numHits := p.findBestSpace()
	// r, c, err := p.getNextSpaceAfter(0, 0)
	// assertOk(err)
	// numHits := -1
	for i := 1; i < 10; i++ {
		if p.numberFits(r, c, i) {
			s := p.Copy()
			s.grid[r][c] = i
			fmt.Println("\\/ trying ", i, " at ", r, ", ", c, " with ", numHits, " hits.")
			s = s.Solve(undos)
			if s != nil {
				return s
			}
		}
	}
	(*undos)++
	fmt.Println("/\\")
	return nil
}

func assertOk(err error) {
	if err != nil {
		panic(err)
	}
}

func printFilenamesInDir() {
	curDir, err := os.Open(".")
	assertOk(err)
	dirNames, err := curDir.Readdirnames(0)
	assertOk(err)
	_, err = fmt.Println(dirNames)
	assertOk(err)
}

func getPuzzleFile() []byte {
	printFilenamesInDir()

	_, err := fmt.Println("Please select a sudoku file.")
	assertOk(err)
	var filename string
	_, err = fmt.Scan(&filename)
	assertOk(err)

	var fileContents []byte
	fileContents, err = ioutil.ReadFile(filename)
	for err != nil {
		_, err = fmt.Println(filename, "didn't work, try again.")
		assertOk(err)
		filename = ""
		fmt.Scanln(&filename)
		assertOk(err)
		fileContents, err = ioutil.ReadFile(filename)
	}
	return fileContents
}

func deleteWhitespace(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\r\n", "", -1)
	s = strings.Replace(s, "+", "", -1)
	s = strings.Replace(s, "|", "", -1)
	s = strings.Replace(s, "-", "", -1)
	return s
}

func loadPuzzle(file []byte) (p *Puzzle, err error) {
	puzzleString := string(file)
	p = NewPuzzle()
	err = p.Load(puzzleString)
	assertOk(err)
	return p, err
}

func main() {
	puzzleFile := getPuzzleFile()
	curPuzzle, err := loadPuzzle(puzzleFile)
	assertOk(err)
	curPuzzle.Print()
	undos := 0
	solvedPuzzle := curPuzzle.Solve(&undos)
	if solvedPuzzle != nil {
		fmt.Println("final after ", undos, " undos")
		solvedPuzzle.Print()
	}
}
