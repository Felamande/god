package main

import (
	"time"

	"github.com/Felamande/god/example/tick"
)

func main() {
	tick.Tick(time.Second * 3)
}
