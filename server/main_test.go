package main

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	loadEnv()
	t.Log(os.Getenv("SSL_ROOT_CERT"))
}
