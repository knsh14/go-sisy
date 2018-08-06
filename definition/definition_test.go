package definition

import "sort"
import "testing"

func TestStandardLibSuccess(t *testing.T) {
	tests := []struct {
		input    [2]string
		expected []string
		comment  string
	}{
		{
			[2]string{"go/ast", "Field"},
			[]string{
				"Doc",
				"Names",
				"Type",
				"Tag",
				"Comment",
			},
			"struct Only Exported Fields",
		},
		{
			[2]string{"net/http/httptest", "ResponseRecorder"},
			[]string{
				"Code",
				"HeaderMap",
				"Body",
				"Flushed",
			},
			"struct contains exported and unexported fields",
		},
		{
			[2]string{"go/token", "FileSet"},
			[]string{},
			"struct only unexported fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.comment, func(test *testing.T) {
			result, err := GetExportedFields(tt.input[0], tt.input[1])
			if err != nil {
				test.Fatal(err)
			}

			sort.Strings(tt.expected)
			sort.Strings(result)
			if len(tt.expected) != len(result) {
				test.Errorf("result fields num is not expected. expected=%d, got=%d", len(tt.expected), len(result))
			}
			for i := range tt.expected {
				if tt.expected[i] != result[i] {
					test.Errorf("unexpected field got. %s", result[i])
				}
			}
		})
	}
}

func TestExternalLibSuccess(t *testing.T) {
	tests := []struct {
		input    [2]string
		expected []string
		comment  string
	}{
		{
			[2]string{"github.com/knsh14/go-struct-initial-fields/testdata", "Foo"},
			[]string{
				"ID",
				"Name",
			},
			"struct Only Exported Fields",
		},
		{
			[2]string{"github.com/knsh14/go-struct-initial-fields/testdata", "Bar"},
			[]string{
				"ID",
				"Age",
			},
			"struct contains exported and unexported fields",
		},
		{
			[2]string{"github.com/knsh14/go-struct-initial-fields/testdata", "Buzz"},
			[]string{},
			"struct only unexported fields",
		},
		{
			[2]string{"../go-struct-initial-fields/testdata", "Buzz"},
			[]string{},
			"struct only unexported fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.comment, func(test *testing.T) {
			result, err := GetExportedFields(tt.input[0], tt.input[1])
			if err != nil {
				test.Fatal(err)
			}

			sort.Strings(tt.expected)
			sort.Strings(result)
			if len(tt.expected) != len(result) {
				test.Errorf("result fields num is not expected. expected=%d, got=%d", len(tt.expected), len(result))
			}
			for i := range tt.expected {
				if tt.expected[i] != result[i] {
					test.Errorf("unexpected field got. %s", result[i])
				}
			}
		})
	}
}

func TestStandardLibFail(t *testing.T) {
	tests := []struct {
		input   [2]string
		comment string
	}{
		{
			[2]string{"go/as", "Field"},
			"wrong package name",
		},
		{
			[2]string{"go/ast", "Fields"},
			"wrong struct name",
		},
		{
			[2]string{"go/ast", "ChanDir"},
			"not struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.comment, func(test *testing.T) {
			_, err := GetExportedFields(tt.input[0], tt.input[1])
			if err == nil {
				test.Fatal("err not found")
			}
			test.Log(err)
		})
	}
}
