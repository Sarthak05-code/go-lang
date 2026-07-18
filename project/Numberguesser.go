package main

import (
	"fmt"
	"math/rand"
	"time"
)

func numberGenerator() int {
	return rand.Intn(100) + 1
}

func gameLogic(difficulty int) {
	fmt.Println("\nLet's start the game!")
	generatedNumber := numberGenerator()
	
	chances := 0
	switch difficulty {
	case 1: 
		chances = 10
	case 2: 
		chances = 5
	case 3: 
		chances = 3
	}

	for i := 0; i < chances; i++ {
		fmt.Printf("\n--- Attempt %d of %d ---\n", i+1, chances)
		fmt.Println("You have 5 seconds for this turn!")
		
		inputChan := make(chan int, 1)
		
		go func() {
			var input int
			for {
				fmt.Print("Enter your guess: ")
				_, err := fmt.Scanln(&input)
				if err != nil || input < 1 || input > 100 {
					fmt.Println("Wrong! Enter a guess between 1 and 100.")
					if err != nil {
						var discard string
						fmt.Scanln(&discard)
					}
					continue
				}
				break
			}
			inputChan <- input
		}()

		var guess int
		timeUp := false

		select {
		case <-time.After(5 * time.Second):
			fmt.Println("\nTime's up! You took too long on this turn.")
			timeUp = true
		case guess = <-inputChan:
		}

		if timeUp {
			continue
		}

		if guess > generatedNumber {
			fmt.Printf("Incorrect! The number is lower than %d\n", guess)
		} else if guess < generatedNumber {
			fmt.Printf("Incorrect! The number is higher than %d\n", guess)
		} else {
			fmt.Printf("Correct! You guessed the correct number in %d attempts!\n", i+1)
			return
		}
	}

	fmt.Printf("\nGame Over! You ran out of chances. The number was %d.\n", generatedNumber)
}

func main() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("You have 5 seconds per turn to guess the number!")

	fmt.Println("\nPlease select the difficulty level: ")
	fmt.Print(`1. Easy (10 chances)
2. Medium (5 chances)
3. Hard (3 chances)
`)

	var choice int
	for {
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scanln(&choice)
		if err != nil || choice < 1 || choice > 3 {
			fmt.Println("Wrong! Enter a choice between 1 and 3.")
			if err != nil {
				var discard string
				fmt.Scanln(&discard)
			}
			continue
		}
		break
	}

	gameLogic(choice)
}