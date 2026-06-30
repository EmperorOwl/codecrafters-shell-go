package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStripBackground(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []string
		wantFields     []string
		wantBackground bool
	}{
		{
			name:           "trailing ampersand",
			tokens:         []string{"sleep", "30", "&"},
			wantFields:     []string{"sleep", "30"},
			wantBackground: true,
		},
		{
			name:           "no background token",
			tokens:         []string{"sleep", "30"},
			wantFields:     []string{"sleep", "30"},
			wantBackground: false,
		},
		{
			name:           "ampersand not last token",
			tokens:         []string{"echo", "&", "hello"},
			wantFields:     []string{"echo", "&", "hello"},
			wantBackground: false,
		},
		{
			name:           "only ampersand",
			tokens:         []string{"&"},
			wantFields:     nil,
			wantBackground: true,
		},
		{
			name:           "empty tokens",
			tokens:         nil,
			wantFields:     nil,
			wantBackground: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, gotBackground := StripBackground(tt.tokens)
			if diff := cmp.Diff(tt.wantBackground, gotBackground); diff != "" {
				t.Errorf("StripBackground() background mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantFields, gotFields, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("StripBackground() fields mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
