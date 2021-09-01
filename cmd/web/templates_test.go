package main

import (
	"testing"
	"time"
)

// TestHumanDate tests that the humanDate function correctly returns a UTC date in our
// human-readable string format. Also, it tests the function returns an empty
// string if time is the zero time value. Also, we test that the function correctly
// converts the time to UTC.
func TestHumanDate(t *testing.T) {
	t.Parallel()
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	// (the want field).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	// Loop over the test cases:
	for _, tt := range tests {
		// rebind tt into this lexical scope to avoid concurrency bug from running
		// sub-tests
		tt := tt

		// use the t.Run() function to run a sub-test for each test case. The first
		// parameter to this is the name of the test (which is used to identify the
		// sub-test in any log output) and the second parameter is an anonymous
		// function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // <== Run each sub-test in parallel
			hd := humanDate(tt.tm)
			t.Logf("testing indexing %q for %q", tt.name, tt.want)
			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
