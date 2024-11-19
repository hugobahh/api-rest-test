package secrets

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadSecrets(sKey string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	environmentPath := filepath.Join(dir, ".local-env")
	err = godotenv.Load(environmentPath)
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return ""
	}

	sTmp := os.Getenv(sKey)
	if sTmp == "" {
		return ""
	} else {
		return sTmp
	}
}
