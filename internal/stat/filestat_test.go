package stat

import "testing"

func TestHumanReadableFilesize(t *testing.T) {
	type test struct {
		input int64
		want  string
	}

	tc := []test{
		{1024, "1.0 kB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.0 PB"},
		{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.0 EB"},
		{1010, "1010 B"},
		{1300, "1.3 kB"},
	}

	for _, tt := range tc {
		got := HumanReadableFilesize(tt.input)
		if got != tt.want {
			t.Errorf("HumanReadableFilesize(%d) = %s, want %s", tt.input, got, tt.want)
		}
	}
}
