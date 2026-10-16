package soch

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetNotesDirName(t *testing.T) {
	tests := []struct {
		name       string
		configured string
		want       string
	}{
		{
			name: "uses default when not configured",
			want: defaultNotesDirName,
		},
		{
			name:       "uses default for whitespace-only value",
			configured: "   ",
			want:       defaultNotesDirName,
		},
		{
			name:       "trims configured value",
			configured: "  project-notes  ",
			want:       "project-notes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resetViper(t)
			if test.configured != "" {
				viper.Set(notesDirNameKey, test.configured)
			}

			if got := getNotesDirName(); got != test.want {
				t.Errorf("getNotesDirName() = %q, want %q", got, test.want)
			}
		})
	}
}
