package main

import "fmt"

// Define a custom enum type
type ValidationError int

// Enum values
const (
	NoError ValidationError = iota
	TooManyEvenNumberError
	TooManyOddNumberError
	TooManyNegativeNumberError
	TooManyPositiveNumberError
	AllSimilarNumberError
)

// Convert enum values to readable strings
func (e ValidationError) String() string {
	switch e {
	case NoError:
		return "No Error"
	case TooManyEvenNumberError:
		return "Too many even numbers."
	case TooManyOddNumberError:
		return "Too many odd numbers."
	case TooManyNegativeNumberError:
		return "Too many negative numbers."
	case TooManyPositiveNumberError:
		return "Too many positive numbers."
	case AllSimilarNumberError:
		return "All the number are the same."
	default:
		return "Unknown Error"
	}
}

func main() {
	var size int

	fmt.Print("Enter the size of your array: ")
	fmt.Scanln(&size)

	array := make([]int, size)

	for i := 0; i < size; i++ {

		for {

			fmt.Printf("Enter number %d: ", i+1)

			_, err := fmt.Scanln(&array[i])

			if err != nil {

				fmt.Println("Invalid input. Please enter an integer.")

				// Discard the invalid input
				var discard string
				fmt.Scanln(&discard)

				continue
			}

			break
		}
	}

	fmt.Printf("\nArray: %v\n", array)

	errors := errorChecker(array)

	if len(errors) == 0 {
		fmt.Println("No errors found.")
	} else {
		fmt.Println("Errors found:")
		for _, err := range errors {
			fmt.Println("-", err)
		}
	}
}

func errorChecker(array []int) []ValidationError {

	var errors []ValidationError

	evenCount := 0
	oddCount := 0
	positiveCount := 0
	negativeCount := 0

	for _, value := range array {

		if value%2 == 0 {
			evenCount++
		} else {
			oddCount++
		}

		if value >= 0 {
			positiveCount++
		} else {
			negativeCount++
		}
	}

	limit := len(array) / 2

	if evenCount > limit {
		errors = append(errors, TooManyEvenNumberError)
	}

	if oddCount > limit {
		errors = append(errors, TooManyOddNumberError)
	}

	if positiveCount > limit {
		errors = append(errors, TooManyPositiveNumberError)
	}

	if negativeCount > limit {
		errors = append(errors, TooManyNegativeNumberError)
	}

	allSame := true

	for i := 1; i < len(array); i++ {
		if array[i] != array[0] {
			allSame = false
			break
		}
	}

	if allSame {
		errors = append(errors, AllSimilarNumberError)
	}

	return errors
}
