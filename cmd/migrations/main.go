package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	direction := os.Args[1]
	switch direction {
	case "up":
		fmt.Println("Running up migrations")
	case "down":
		fmt.Println("Running down migrations")
	default:
		log.Fatal("Invalid direction, must be up or down")
	}

}
