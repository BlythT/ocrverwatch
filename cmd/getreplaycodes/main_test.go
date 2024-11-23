package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetReplayCodes(t *testing.T) {
	// Define test cases
	tests := []struct {
		filePath      string
		expectedCodes []string
	}{
		{"assets/test/screenshot1.jpg", []string{"9MKHMB", "859MR0", "QGMZ0K", "6MZ63A", "9NAVY1", "9NRJ59"}},
		{"assets/test/screenshot2.png", []string{"2HYX1G", "9WHDDS", "0Z9X7N", "E9E4H1", "N8NNXS", "TJK3X4"}},
		{"assets/test/screenshot3.png", []string{"72JGBF", "7TSX7W", "FHXC4Y"}},
		{"assets/test/photo1.png", []string{"2P1C4D", "1WJXEG", "WCXJQY", "Z3X710", "1PN311", "SRSRE3"}},
		{"assets/test/photo2.jpg", []string{"YTVFW1", "OREOXM", "ZE299H", "CDE1AW", "38Y4CT", "4HN3NX", "HZX6BJ"}},

	}

	// Iterate over test cases
	for _, test := range tests {
		t.Run(filepath.Base(test.filePath), func(t *testing.T) {
			// Ensure the file exists before testing
			if _, err := os.ReadFile(test.filePath); err != nil {
				t.Fatalf("failed to read file %s: %v", test.filePath, err)
			}

			// Call the function to be tested
			replayCodes, err := getReplayCodes(test.filePath)
			if err != nil {
				t.Fatalf("GetReplayCodes returned an error: %v", err)
			}

			for _, code := range replayCodes {
				fmt.Println(code)
			}

			// Check if the returned codes match the expected codes
			if len(replayCodes) != len(test.expectedCodes) {
				t.Fatalf("expected %d codes, got %d (%v)", len(test.expectedCodes), len(replayCodes), replayCodes)
			}

			for i, code := range replayCodes {
				if code != test.expectedCodes[i] {
					t.Errorf("expected code %q, got %q", test.expectedCodes[i], code)
				}
			}
		})
	}
}
