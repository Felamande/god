package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello dw")
	id := 0
	for {
		fmt.Println("tick", id+1)
		time.Sleep(time.Second)
		id++

	}

}
