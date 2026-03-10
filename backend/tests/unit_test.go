package tests

import (
	"testing"
	"time"

	"github.com/pizzattolog/licencas/internal/auth"
	"github.com/pizzattolog/licencas/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── JWT Tests ───────────────────────────────────────────────────────────────

func TestJWT_GenerateAndValidate(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-com-32-chars-minimo!!")
	defer os.Unsetenv("JWT_SECRET")

	svc := auth.NewJWTService()

	token, err := svc.GenerateToken(1, "teste@pizzattolog.com.br", "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "teste@pizzattolog.com.br", claims.Email)
	assert.Equal(t, "admin", claims.Role)
}

func TestJWT_TokenInvalido(t *testing.T) {
	svc := auth.NewJWTService()
	_, err := svc.ValidateToken("token-completamente-invalido")
	assert.Error(t, err)
}

// ─── Model Tests ─────────────────────────────────────────────────────────────

func TestUsuario_SetVerificarSenha(t *testing.T) {
	u := &models.Usuario{}
	err := u.SetSenha("MinhaSenha@123")
	require.NoError(t, err)
	assert.NotEmpty(t, u.SenhaHash)
	assert.NotEqual(t, "MinhaSenha@123", u.SenhaHash)

	assert.True(t, u.VerificarSenha("MinhaSenha@123"))
	assert.False(t, u.VerificarSenha("SenhaErrada"))
}

func TestLicenca_AtualizarStatus(t *testing.T) {
	tests := []struct {
		name     string
		dias     int
		expected models.StatusLicenca
	}{
		{"vencida", -1, models.StatusVencida},
		{"critica 7 dias", 7, models.StatusProximaVencimento},
		{"proxima 30 dias", 30, models.StatusProximaVencimento},
		{"ativa 60 dias", 60, models.StatusAtiva},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &models.Licenca{
				DataValidade: time.Now().AddDate(0, 0, tt.dias),
			}
			l.AtualizarStatus()
			assert.Equal(t, tt.expected, l.Status)
		})
	}
}

func TestLicenca_DiasParaVencer(t *testing.T) {
	l := &models.Licenca{
		DataValidade: time.Now().AddDate(0, 0, 15),
	}
	dias := l.DiasParaVencer()
	// Tolerância de 1 dia por conta de frações de hora
	assert.True(t, dias >= 14 && dias <= 15, "esperava 14 ou 15, obteve %d", dias)
}
