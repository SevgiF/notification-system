package defer_util

import "log"

func DeferWithErrorHandling(f func() error) {
	defer func() {
		// panic durumlarını yakalar
		if err := recover(); err != nil {
			log.Printf("deferred function panic error: %v", err)
		}
	}()

	if err := f(); err != nil {
		log.Printf("deferred function error: %v", err)
	}
}
