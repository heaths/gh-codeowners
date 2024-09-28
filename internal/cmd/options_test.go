package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNumberRef(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		number  string
		wantErr bool
		want    int
	}{
		{
			name:   "number",
			number: "1234",
			want:   1234,
		},
		{
			name:   "#number",
			number: "#1234",
			want:   1234,
		},
		{
			name:    "invalid",
			number:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNumberRef(tt.number)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
