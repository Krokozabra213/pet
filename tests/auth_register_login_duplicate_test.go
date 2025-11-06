package tests

import (
	"testing"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)
	st.Cleanup(func() {
		st.CleanupTestData()
	})
	appID, err := st.CreateApp("test")
	require.NoError(t, err)
	require.NotEmpty(t, appID)

	respPublicKey, err := st.AuthClient.GetPublicKey(ctx, &sso.PublicKeyRequest{
		AppId: int32(appID),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respPublicKey.GetPublicKey())

	publicKeyManager, err := keymanager.NewPublic([]byte(respPublicKey.GetPublicKey()))
	require.NoError(t, err)
	assert.NotEmpty(t, publicKeyManager)

	username := randomUsername()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}
