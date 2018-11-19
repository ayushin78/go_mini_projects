package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

var fileName = "problems.csv"

func main() {

	var userAnswer string
	var score int
	var timeLimit = 3 * time.Second

	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Failed to open the csv file")
		os.Exit(1)
	}

	defer csvFile.Close()

	records, _ := csv.NewReader(csvFile).ReadAll()
	problems := parseRecords(records)

	ansChan := make(chan string)

	for _, p := range problems {
		timer := time.NewTimer(timeLimit)

		go askQuestion(p, ansChan)

		select {
		case <-timer.C:
			fmt.Println("TIME UP")
			break
		case userAnswer = <-ansChan:
			score = checkAns(p.answer, userAnswer, score)
		}
	}

	fmt.Printf("You scored %v out of %v\n", score, len(problems))
}

func checkAns(ans string, userAnswer string, score int) int {
	if strings.Compare(ans, userAnswer) == 0 {
		score = score + 1
		fmt.Printf("Correct Answer! And your score is %v \n", score)
	} else {
		fmt.Printf("Wrong answer! Correct answer is %v \n", ans)
	}
	return score
}

func askQuestion(p problem, ansChan chan<- string) {
	var userAnswer string

	fmt.Printf("what is %v?\n", p.question)
	fmt.Scanln(&userAnswer)

	ansChan <- userAnswer
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{
			question: record[0],
			answer:   record[1],
		}
	}
	return problems
}
