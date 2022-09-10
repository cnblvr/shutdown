package sleep

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func Do(duration time.Duration) {
	if duration <= 0 {
		return
	}
	id := fmt.Sprintf("%08x", rand.Uint32())
	log.Printf("[%s] start sleep %s", id, duration)
	time.Sleep(duration)
	log.Printf("[%s] finish sleep %s", id, duration)
}
