package ts

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	timestampTimer sync.Once
	timestamp      uint32
)

// Timestamp returns the current time.
// Make sure to start the updater once using StartTimeStampUpdater() before calling.
func Timestamp() uint32 {
	return atomic.LoadUint32(&timestamp)
}

// StartTimeStampUpdater starts a concurrent function which stores the timestamp to an atomic value per second,
// which is much better for performance than determining it at runtime each time
func StartTimeStampUpdater() {
	timestampTimer.Do(func() {
		// set initial value
		atomic.StoreUint32(&timestamp, uint32(time.Now().Unix()))
		go func(sleep time.Duration) {
			ticker := time.NewTicker(sleep)
			defer ticker.Stop()

			for t := range ticker.C {
				// update timestamp
				atomic.StoreUint32(&timestamp, uint32(t.Unix()))
			}
		}(1 * time.Second) // duration
	})
}
