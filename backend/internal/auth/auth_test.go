// backend/internal/auth/auth_test.go
package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("mypassword")
	assert.NoError(t, err)
	assert.NotEqual(t, "mypassword", hash)

	assert.True(t, VerifyPassword(hash, "mypassword"))
	assert.False(t, VerifyPassword(hash, "wrong"))
}

func TestJWTIssueAndParse(t *testing.T) {
	issuer := NewJWTIssuer("super-secret-key")
	token, err := issuer.Issue(42, "super_admin", time.Hour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := issuer.Parse(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), claims.UID)
	assert.Equal(t, "super_admin", claims.Role)
}

func TestJWTParse_Invalid(t *testing.T) {
	issuer := NewJWTIssuer("super-secret-key")
	_, err := issuer.Parse("not-a-token")
	assert.Error(t, err)
}

func TestJWTParse_WrongSecret(t *testing.T) {
	issuer1 := NewJWTIssuer("secret-a")
	token, _ := issuer1.Issue(1, "admin", time.Hour)

	issuer2 := NewJWTIssuer("secret-b")
	_, err := issuer2.Parse(token)
	assert.Error(t, err)
}
