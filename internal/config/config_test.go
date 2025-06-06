package config

import (
	"testing"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/pagu/internal/platforms/discord"
	"github.com/pagu-project/pagu/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

// TestBasicCheck tests the BasicCheck method of the Config struct.
func TestBasicCheck(t *testing.T) {
	// Create a temporary directory for the WalletPath
	tempWalletPath := t.TempDir()

	// Define test cases
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid config",
			cfg: Config{
				Engine: engine.Config{
					Wallet: wallet.Config{
						Address:  "test_wallet_address",
						Path:     tempWalletPath, // Use the temporary directory
						Password: "test_password",
					},
					NetworkNodes: []string{"http://127.0.0.1:8545"},

					Phoenix: phoenix.Config{},
				},
				Discord: discord.Config{
					Token:   "MTEabc123",
					GuildID: "123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "No RPCNodes",
			cfg: Config{
				Engine: engine.Config{
					Wallet: wallet.Config{
						Address:  "test_wallet_address",
						Path:     "/valid/path",
						Password: "test_password",
					},
					NetworkNodes: []string{},
				},
				Discord: discord.Config{
					Token:   "MTEabc123",
					GuildID: "123456789",
				},
			},
			wantErr: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Perform the check
			err := tt.cfg.BasicCheck()

			// Assert the error based on wantErr
			if tt.wantErr {
				assert.Error(t, err, "Config.BasicCheck() should return an error")
			} else {
				assert.NoError(t, err, "Config.BasicCheck() should not return an error")
			}
		})
	}
}
