package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

// getEnvOrFail, belirtilen ortam değişkenini alır, yoksa hata verir.
func GetEnvOrFail(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Ortam değişkeni eksik: %s", key)
	}
	return value
}

// getIntEnv, ortam değişkenini int türünde alır veya varsayılan değeri döner.
func GetIntEnv(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Ortam değişkeni '%s' geçersiz, varsayılan değer kullanılacak: %d", key, defaultValue)
	}
	return defaultValue
}

// getDurationEnv, ortam değişkenini time.Duration türünde alır veya varsayılan değeri döner.
func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Ortam değişkeni '%s' geçersiz, varsayılan değer kullanılacak: %s", key, defaultValue)
	}
	return defaultValue
}
