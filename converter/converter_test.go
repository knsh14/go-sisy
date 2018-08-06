package converter

import (
	"go/format"
	"os"
	"testing"
)

func TestSample(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "../testdata/testdata.go",
		},
	}

	for _, tt := range tests {
		f, fset, err := Convert(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		format.Node(os.Stdout, fset, f)
	}
}
