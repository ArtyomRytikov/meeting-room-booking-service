package auth

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken("user")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if claims.Role != "user" {
		t.Fatalf("expected role=user, got %s", claims.Role)
	}

	if claims.UserID == "" {
		t.Fatal("expected non-empty user id")
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestFixedUserIDByRole(t *testing.T) {
	id, err := FixedUserIDByRole("admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty admin id")
	}
}
