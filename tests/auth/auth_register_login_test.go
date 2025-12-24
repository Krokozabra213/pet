package tests

import (
	"testing"
	"time"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"github.com/Krokozabra213/sso/tests/auth/suite"
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

	assert.Equal(t, respPublicKey.GetPublicKey(), st.PublicPEM)

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
	accessExp := loginTime.Add(st.Cfg.Auth.JWT.AccessTokenTTL)
	refreshExp := loginTime.Add(st.Cfg.Auth.JWT.RefreshTokenTTL)

	accessToken := respLogin.GetAccessToken()
	refreshToken := respLogin.GetRefreshToken()
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	accessData, err := st.JWTmanager.ParseAccess(accessToken)
	require.NoError(t, err)
	assert.NotNil(t, accessData)
	assert.Equal(t, respReg.GetUserId(), int64(accessData.UserID))
	assert.Equal(t, appID, accessData.AppID)
	assert.Equal(t, username, accessData.Username)
	assert.WithinDuration(t, accessExp, accessData.Exp, 10*time.Second)

	refreshData, err := st.JWTmanager.ParseRefresh(refreshToken)
	require.NoError(t, err)
	assert.NotNil(t, refreshData)
	assert.Equal(t, respReg.GetUserId(), int64(refreshData.UserID))
	assert.Equal(t, appID, refreshData.AppID)
	assert.Equal(t, username, refreshData.Username)
	assert.NotEmpty(t, refreshData.JWTID)
	assert.WithinDuration(t, refreshExp, refreshData.Exp, 10*time.Second)

	// access token завершился, сервер обращается за новой парой
	respRefresh, err := st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)

	refreshTime := time.Now()
	accessExp = refreshTime.Add(st.Cfg.Auth.JWT.AccessTokenTTL)
	refreshExp = refreshTime.Add(st.Cfg.Auth.JWT.RefreshTokenTTL)

	accessToken = respRefresh.GetAccessToken()
	refreshToken = respRefresh.GetRefreshToken()
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	accessData, err = st.JWTmanager.ParseAccess(accessToken)
	require.NoError(t, err)
	assert.NotNil(t, accessData)
	assert.Equal(t, respReg.GetUserId(), int64(accessData.UserID))
	assert.Equal(t, appID, accessData.AppID)
	assert.Equal(t, username, accessData.Username)
	assert.WithinDuration(t, accessExp, accessData.Exp, 10*time.Second)

	refreshData, err = st.JWTmanager.ParseRefresh(refreshToken)
	require.NoError(t, err)
	assert.NotNil(t, refreshData)
	assert.Equal(t, respReg.GetUserId(), int64(refreshData.UserID))
	assert.Equal(t, appID, refreshData.AppID)
	assert.Equal(t, username, refreshData.Username)
	assert.NotEmpty(t, refreshData.JWTID)
	assert.WithinDuration(t, refreshExp, refreshData.Exp, 10*time.Second)

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
