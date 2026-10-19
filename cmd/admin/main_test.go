package main

import (
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Change to project root so configs/application.yaml can be found
	err := os.Chdir("../../")
	if err != nil {
		t.Skipf("cannot change directory: %v", err)
	}

	// Since main() calls log.Fatalf on config error, we test that loadConfig works.
	cfg, err := loadConfig("configs/application.yaml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Database.Name != "codeflix_catalog" {
		t.Errorf("expected database codeflix_catalog, got %s", cfg.Database.Name)
	}
}
