package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntSlice []int

const asize = 256

/******** utility functions ********/

func readFile(fileName string) []string {

	data := make([]string, 0)

	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(fi)

	scanner := bufio.NewScanner(r)

	temp := make([]string, asize)
	counter := 0

	for scanner.Scan() {
		in := scanner.Text()
		counter++
		temp[counter-1] = in

		if counter == asize {
			counter = 0
			data = append(data, temp...)
		}
	}

	temp = temp[:counter]
	data = append(data, temp...)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading problem:", err)
	}

	return data
}

func parseDimacs(input []string) [][]int {
	counter := 0
	data := make([][]int, len(input))
	for i := range input {
		input[i] = strings.TrimSpace(input[i])
		c := string(input[i][0])
		if c != "p" && c != "c" && c != "0" && c != "%" {
			counter++
			s := strings.Split(input[i], " ")
			s = s[:len(s)-1]
			data[counter-1] = make([]int, len(s))

			for j := range s {
				data[counter-1][j], _ = strconv.Atoi(s[j])
			}
		}

	}
	data = data[:counter]
	return data
}

func contains(xs []int, l int) bool {
	isin := false
	for i := range xs {
		if xs[i] == l {
			isin = true
		}
	}
	return isin
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

/******** dpll ********/

func simplify(f [][]int, l int) [][]int {
	var simplifyClause = func(c []int, l int) []int {
		counter := 0
		out := make([]int, len(c))
		for i := range c {
			if c[i] != -l {
				out[counter] = c[i]
				counter++
			}
		}
		return out[:counter]
	}

	nf := make([][]int, len(f))
	counter := 0
	for i := range f {
		if !contains(f[i], l) {
			nc := simplifyClause(f[i], l)
			nf[counter] = nc
			counter++
		}
	}
	return nf[:counter]
}

func chooseLiteral(f [][]int) int {
	out := 0
	for i := range f {
		if len(f[i]) > 0 {
			out = f[i][0]
			break
		}
	}
	if out == 0 {
		panic("No literals")
	}
	return out
}

func unitpropagate(f [][]int, r []int) ([][]int, []int) {
	u := 0
	for i := range f {
		if len(f[i]) == 1 {
			u = f[i][0]
			break
		}
	}
	if Abs(u) > 0 {
		return unitpropagate(simplify(f, u), append(r, u))
	} else {
		return f, r
	}
}

/****
TODO
func pureLiteralAssign(f [][]int, r []int) ([][]int, []int) {

}
****/

func dpll(f [][]int, r []int) []int {
	var containsEmpty = func(f [][]int) bool {
		cns := false
		for i := range f {
			if len(f[i]) == 0 {
				cns = true
				break
			}
		}
		return cns
	}

	f, r = unitpropagate(f, r)

	if len(f) == 0 {
		return r
	} else if containsEmpty(f) {
		return []int{}
	} else {
		l := chooseLiteral(f)
		t := dpll(simplify(f, l), append(r, l))
		if len(t) != 0 {
			return t
		} else {
			return dpll(simplify(f, -l), append(r, -l))
		}
	}
}

func main() {

	if len(os.Args) > 1 {
		rf := readFile(os.Args[1])
		formula := parseDimacs(rf)

		rj := dpll(formula, []int{})

		fmt.Print("v")
		for i := range rj {
			fmt.Print(" ")
			fmt.Print(rj[i])
		}
	}
}
