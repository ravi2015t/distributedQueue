package web

import "testing"

func TestInvalidCategory(t *testing.T) {
	testCases := []struct {
		category string
		valid    bool
	}{
		{category: "", valid: false},
		{category: ".", valid: false},
		{category: "..", valid: false},
		{category: "numbers", valid: true},
		{category: ":_numbers", valid: true},
	}

	for _, tc := range testCases {
		got := isValidCategory(tc.category)
		want := tc.valid

		if got != want {
			t.Errorf("isValidCategory(%q) = %v; want %v", tc.category, got, want)
		}
	}
}
