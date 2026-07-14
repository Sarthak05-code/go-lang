package main

import (
	"fmt"
	"time"
)

func counter(stopChan chan bool) {
	i := 0
	for {
		select {
		case <-stopChan:
			// If we receive a signal on the channel, stop the loop
			fmt.Printf("\nThe program ran for %d seconds\n", i)
			return
		default:
			// If no signal, wait 1 second, then increment
			time.Sleep(1 * time.Second)
			i++
		}
	}
}


func main() {
	// Create a channel to communicate between main and the goroutine
	stopChan := make(chan bool)

	// Pass the channel to the goroutine
	go counter(stopChan)

	var value string
	fmt.Println("Enter name: ")
	fmt.Scanln(&value)

	// Send a signal to the goroutine to tell it to stop
	stopChan <- true

	fmt.Printf("Your name is: %s\n", value)
}