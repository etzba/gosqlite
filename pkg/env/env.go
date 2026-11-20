package env

import "os"

var ReleaseStage = GetEnv("", "RELEASE_STAGE")

func GetEnv(defaultVal, envStr string) string {
	if os.Getenv(envStr) != "" {
		return os.Getenv(envStr)
	}
	return defaultVal
}
