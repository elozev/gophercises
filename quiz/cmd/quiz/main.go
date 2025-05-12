package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	filename = flag.String("filename", "problems.csv", "csv file with problems")
	limit = flag.Int("limit", 10, "amount of time in seconds allowed to solve the problems")
)

// func getUserInput(question string) string {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Printf("%s?\n>",question)
// 	scanner.Scan()
// 	return scanner.Text()
// }

func getUserInput(q string) string {
	fmt.Printf("%s?\n>", q)
	var answer string
	fmt.Scanf("%s\n", &answer)

	return answer
}

type Problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []Problem {
	var rat = make([]Problem, len(lines))

	for i, line := range lines {
		rat[i] = Problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return rat
}

func parseCSV() {
	file, err := os.Open(*filename)
	defer file.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "something went opening %s\n", filename)
		os.Exit(-1)
	}
	
	reader := csv.NewReader(file)

	var correctAnswers int

	lines, err := reader.ReadAll()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing csv lines", err)
	}

	problems := parseLines(lines)

	defer func() {
		fmt.Printf("You scored %d out of %d\n", correctAnswers, len(problems))
	}()

	for _, p := range problems {
		userInput := getUserInput(p.q)

		if userInput == p.a {
			correctAnswers++
		}
	}
}

func main() {
	flag.Parse()

	parseCSV()
}