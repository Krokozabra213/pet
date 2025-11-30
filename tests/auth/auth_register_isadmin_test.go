package tests

import (
	"testing"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	"github.com/Krokozabra213/sso/tests/auth/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin_IsAdmin(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	username := randomUsername()
	pass := randomFakePassword()

	// регистрация пользователя
	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	// добавляем пользователя в админы
	adminID, err := st.CreateAdmin(respReg.GetUserId())
	require.NoError(t, err)
	assert.NotEmpty(t, adminID)

	// проверяем что он стал админом
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &sso.IsAdminRequest{
		UserId: adminID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respIsAdmin.IsAdmin)
	assert.Equal(t, respIsAdmin.IsAdmin, true)
}
