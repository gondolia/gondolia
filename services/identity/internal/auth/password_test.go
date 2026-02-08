package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "admin123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-that-should-still-work-fine",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
			if !tt.wantErr && hash == tt.password {
				t.Error("HashPassword() returned unhashed password")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// Pre-hash passwords for testing
	adminHash, _ := HashPassword("admin123")
	testHash, _ := HashPassword("test123")

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct admin password",
			password: "admin123",
			hash:     adminHash,
			want:     true,
		},
		{
			name:     "correct test password",
			password: "test123",
			hash:     testHash,
			want:     true,
		},
		{
			name:     "wrong password",
			password: "wrongpassword",
			hash:     adminHash,
			want:     false,
		},
		{
			name:     "empty password against hash",
			password: "",
			hash:     adminHash,
			want:     false,
		},
		{
			name:     "password against invalid hash",
			password: "admin123",
			hash:     "invalidhash",
			want:     false,
		},
		{
			name:     "password against empty hash",
			password: "admin123",
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyPassword(tt.password, tt.hash); got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password 8 chars",
			password: "admin123",
			wantErr:  false,
		},
		{
			name:     "valid long password",
			password: "this-is-a-secure-password",
			wantErr:  false,
		},
		{
			name:     "too short password",
			password: "short",
			wantErr:  true,
		},
		{
			name:     "7 char password",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSecureToken(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "16 bytes", length: 16},
		{name: "32 bytes", length: 32},
		{name: "64 bytes", length: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateSecureToken(tt.length)
			if err != nil {
				t.Errorf("GenerateSecureToken() error = %v", err)
				return
			}

			// Hex encoding doubles the length
			expectedLen := tt.length * 2
			if len(token) != expectedLen {
				t.Errorf("GenerateSecureToken() length = %d, want %d", len(token), expectedLen)
			}

			// Verify uniqueness
			token2, _ := GenerateSecureToken(tt.length)
			if token == token2 {
				t.Error("GenerateSecureToken() generated duplicate tokens")
			}
		})
	}
}

func TestHashToken(t *testing.T) {
	token := "test-token-12345"
	hash1 := HashToken(token)
	hash2 := HashToken(token)

	// Same token should produce same hash
	if hash1 != hash2 {
		t.Error("HashToken() not deterministic")
	}

	// Different tokens should produce different hashes
	hash3 := HashToken("different-token")
	if hash1 == hash3 {
		t.Error("HashToken() collision for different tokens")
	}

	// Hash should be 64 chars (SHA256 = 32 bytes = 64 hex chars)
	if len(hash1) != 64 {
		t.Errorf("HashToken() length = %d, want 64", len(hash1))
	}
}
