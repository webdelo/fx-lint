package plugin

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestNew_NilSettings_Errors(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil settings")
	}
}

func TestNew_EmptyModPath_Errors(t *testing.T) {
	_, err := New(map[string]any{"mod_path": ""})
	if err == nil {
		t.Fatal("expected error for empty mod_path")
	}
}

func TestNew_Valid_BuildsThreeAnalyzers(t *testing.T) {
	p, err := New(map[string]any{"mod_path": "example.com/proj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	as, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers: %v", err)
	}
	if got, want := len(as), 3; got != want {
		t.Fatalf("analyzer count = %d, want %d", got, want)
	}
}

func TestPlugin_LoadMode(t *testing.T) {
	p, err := New(map[string]any{"mod_path": "x/y"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.GetLoadMode(), register.LoadModeTypesInfo; got != want {
		t.Fatalf("load mode = %q, want %q", got, want)
	}
}
