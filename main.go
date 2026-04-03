package main

import (
	"fmt"
)

// I'd initialize this on a separate file
// But I don't know how you'd do it in go, so let's just keep this here.
type Task struct {
	ID     int
	Title  string
	Status string
}

func main() {
	task1 := Task{
		ID:     1,
		Title:  "Learn Go",
		Status: "in-progress",
	}

	fmt.Println("Hello", task1.ID)
	fmt.Println("Hello", task1.Title)
	fmt.Println("Hello", task1.Status)
}
