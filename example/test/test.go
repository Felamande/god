package tick

import (
	"fmt"
	"time"
)

func Tick(d time.Duration) {
	id := 0
	for {
		fmt.Println("tick", id+4)
		time.Sleep(d)
		id++
	}
}
