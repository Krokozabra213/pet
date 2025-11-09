package tests

import (
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/proto/sso"
	"github.com/Krokozabra213/sso/internal/auth/lib/jwt"
	keymanager "github.com/Krokozabra213/sso/internal/auth/lib/key-manager"
	"github.com/Krokozabra213/sso/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	// регистрируем приложение которое будет пользоваться нашим сервисом авторизации
	appID, err := st.CreateApp("test")
	require.NoError(t, err)
	assert.NotEmpty(t, appID)

	// получаем публичный ключ rsa для парсинга jwt токенов
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

	// регистрация пользователя
	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	// логин пользователя
	respLogin, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: username,
		Password: pass,
		AppId:    int32(appID),
	})
	require.NoError(t, err)

	loginTime := time.Now()
	accessExp := loginTime.Add(time.Duration(st.Cfg.Security.AccessTokenTTL) * time.Second)
	refreshExp := loginTime.Add(time.Duration(st.Cfg.Security.RefreshTokenTTL) * time.Second)

	accessToken := respLogin.GetAccessToken()
	refreshToken := respLogin.GetRefreshToken()
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	accessParsedData, err := jwt.ParseAccess(accessToken, publicKeyManager)
	require.NoError(t, err)
	assert.NotNil(t, accessParsedData)
	assert.Equal(t, respReg.GetUserId(), int64(accessParsedData.UserID))
	assert.Equal(t, appID, accessParsedData.AppID)
	assert.Equal(t, username, accessParsedData.Username)
	assert.WithinDuration(t, accessExp, accessParsedData.Exp, 10*time.Second)

	refreshParsedData, err := jwt.ParseRefresh(refreshToken, publicKeyManager)
	require.NoError(t, err)
	assert.NotNil(t, refreshParsedData)
	assert.Equal(t, respReg.GetUserId(), int64(refreshParsedData.UserID))
	assert.Equal(t, appID, refreshParsedData.AppID)
	assert.Equal(t, username, refreshParsedData.Username)
	assert.NotEmpty(t, refreshParsedData.JWTID)
	assert.WithinDuration(t, refreshExp, refreshParsedData.Exp, 10*time.Second)

	// access token завершился, сервер обращается за новой парой
	respRefresh, err := st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)

	refreshTime := time.Now()
	accessExp = refreshTime.Add(time.Duration(st.Cfg.Security.AccessTokenTTL) * time.Second)
	refreshExp = refreshTime.Add(time.Duration(st.Cfg.Security.RefreshTokenTTL) * time.Second)

	accessToken = respRefresh.GetAccessToken()
	refreshToken = respRefresh.GetRefreshToken()
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	accessParsedData, err = jwt.ParseAccess(accessToken, publicKeyManager)
	require.NoError(t, err)
	assert.NotNil(t, accessParsedData)
	assert.Equal(t, respReg.GetUserId(), int64(accessParsedData.UserID))
	assert.Equal(t, appID, accessParsedData.AppID)
	assert.Equal(t, username, accessParsedData.Username)
	assert.WithinDuration(t, accessExp, accessParsedData.Exp, 10*time.Second)

	refreshParsedData, err = jwt.ParseRefresh(refreshToken, publicKeyManager)
	require.NoError(t, err)
	assert.NotNil(t, refreshParsedData)
	assert.Equal(t, respReg.GetUserId(), int64(refreshParsedData.UserID))
	assert.Equal(t, appID, refreshParsedData.AppID)
	assert.Equal(t, username, refreshParsedData.Username)
	assert.NotEmpty(t, refreshParsedData.JWTID)
	assert.WithinDuration(t, refreshExp, refreshParsedData.Exp, 10*time.Second)

	// пользователь выходит с сайта
	respLogout, err := st.AuthClient.Logout(ctx, &sso.LogoutRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respLogout.GetSuccess())
	assert.Equal(t, respLogout.GetSuccess(), true)

	//проверяем что токен отозван
	respRefresh, err = st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: refreshToken,
	})
	require.Error(t, err)
	assert.ErrorContains(t, err, "token revoked")
	assert.Empty(t, respRefresh)
}
