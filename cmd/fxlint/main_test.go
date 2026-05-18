package main

import (
	"testing"
)

func TestExtractModPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		fallback  string
		wantValue string
		wantRest  []string
	}{
		{
			name:      "nil args empty fallback",
			args:      nil,
			fallback:  "",
			wantValue: "",
			wantRest:  []string{},
		},
		{
			name:      "empty args with fallback",
			args:      []string{},
			fallback:  "x",
			wantValue: "x",
			wantRest:  []string{},
		},
		{
			name:      "dash-mod-path separate value",
			args:      []string{"-mod-path", "x/y", "./..."},
			fallback:  "",
			wantValue: "x/y",
			wantRest:  []string{"./..."},
		},
		{
			name:      "dash-mod-path equals form",
			args:      []string{"-mod-path=x/y", "./..."},
			fallback:  "",
			wantValue: "x/y",
			wantRest:  []string{"./..."},
		},
		{
			name:      "double-dash equals form",
			args:      []string{"--mod-path=x/y", "./..."},
			fallback:  "",
			wantValue: "x/y",
			wantRest:  []string{"./..."},
		},
		{
			name:      "passthrough flags preserved",
			args:      []string{"-flags", "-mod-path", "x"},
			fallback:  "",
			wantValue: "x",
			wantRest:  []string{"-flags"},
		},
		{
			name:      "last mod-path wins",
			args:      []string{"-mod-path", "x", "-mod-path=y"},
			fallback:  "",
			wantValue: "y",
			wantRest:  []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotValue, gotRest := extractModPath(tc.args, tc.fallback)
			if gotValue != tc.wantValue {
				t.Errorf("value = %q, want %q", gotValue, tc.wantValue)
			}
			if len(gotRest) != len(tc.wantRest) {
				t.Errorf("rest = %v, want %v", gotRest, tc.wantRest)
				return
			}
			for i := range gotRest {
				if gotRest[i] != tc.wantRest[i] {
					t.Errorf("rest[%d] = %q, want %q", i, gotRest[i], tc.wantRest[i])
				}
			}
		})
	}
}
