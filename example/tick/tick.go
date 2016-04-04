package tick

import (
	"fmt"
	"time"
)

func Tick(d time.Duration) {
	id := 0
	for {
		fmt.Printf("tick every %v %d\n", d, id)
		time.Sleep(d)
		id += 5
	}
}
