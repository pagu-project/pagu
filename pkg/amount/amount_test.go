package amount_test

import (
	"encoding/json"
	"testing"

	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected amount.Amount
		err      bool
	}{
		{
			name:     "Valid string input",
			data:     []byte(`"123.456"`),
			expected: amount.Amount(123456000000),
			err:      false,
		},
		{
			name:     "Valid float input",
			data:     []byte(`123.456`),
			expected: amount.Amount(123456000000),
			err:      false,
		},
		{
			name:     "Invalid string input",
			data:     []byte(`"invalid"`),
			expected: amount.Amount(0),
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var amt amount.Amount
			err := yaml.Unmarshal(tt.data, &amt)

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, amt)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected amount.Amount
		err      bool
	}{
		{
			name:     "Valid string input",
			data:     []byte(`"123.456"`),
			expected: amount.Amount(123456000000),
			err:      false,
		},
		{
			name:     "Valid float input",
			data:     []byte(`123.456`),
			expected: amount.Amount(123456000000),
			err:      false,
		},
		{
			name:     "Invalid string input",
			data:     []byte(`"invalid"`),
			expected: amount.Amount(0),
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var amt amount.Amount
			err := json.Unmarshal(tt.data, &amt)

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, amt)
			}
		})
	}
}
