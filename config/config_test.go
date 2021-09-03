package config

import (
	"fmt"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	name := "cache.json"

	SetupConfig(name)

	fmt.Println(GlobalConfig)
}
