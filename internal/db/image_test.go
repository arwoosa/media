package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithImageVariants(t *testing.T) {
	tests := []struct {
		name     string
		variants []string
		expected map[string]string
	}{
		{
			name: "Valid variants",
			variants: []string{
				"https://imagedelivery.net/2_2N5I_Gk2ffs-N2-w/e2f9e558-5364-4465-9329-ea1b5f185900/public",
				"https://imagedelivery.net/2_2N5I_Gk2ffs-N2-w/e2f9e558-5364-4465-9329-ea1b5f185900/thumbnail",
			},
			expected: map[string]string{
				"public":    "/cdn-images/e2f9e558-5364-4465-9329-ea1b5f185900/public",
				"thumbnail": "/cdn-images/e2f9e558-5364-4465-9329-ea1b5f185900/thumbnail",
			},
		},
		{
			name:     "Empty variants",
			variants: []string{},
			expected: map[string]string{},
		},
		{
			name: "Variants with invalid URL",
			variants: []string{
				"https://imagedelivery.net/2_2N5I_Gk2ffs-N2-w/e2f9e558-5364-4465-9329-ea1b5f185900/public",
				"://invalid-url",
			},
			expected: map[string]string{
				"public": "/cdn-images/e2f9e558-5364-4465-9329-ea1b5f185900/public",
			},
		},
		{
			name: "URL with no path",
			variants: []string{
				"https://example.com",
			},
			expected: map[string]string{},
		},
		{
			name: "URL with trailing slash",
			variants: []string{
				"https://imagedelivery.net/2_2N5I_Gk2ffs-N2-w/e2f9e558-5364-4465-9329-ea1b5f185900/public/",
			},
			expected: map[string]string{
				"public": "/cdn-images/e2f9e558-5364-4465-9329-ea1b5f185900/public/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := NewImage(WithImageVariants(tt.variants))
			assert.Equal(t, tt.expected, img.Variants)
		})
	}
}
