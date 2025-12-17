package tests

import (
	"log"
	"testing"

	"github.com/Krokozabra213/protos/gen/go/sso"
	authBusiness "github.com/Krokozabra213/sso/internal/auth/business"
	"github.com/Krokozabra213/sso/tests/auth/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin_Errors(t *testing.T) {
	ctx, st := suite.New(t)
	t.Cleanup(func() {
		st.CleanupTestData()
	})

	appID := randomID(1, 100_000)

	// проверка getpublickey на ошибки
	respPublicKey, err := st.AuthClient.GetPublicKey(ctx, &sso.PublicKeyRequest{
		AppId: int32(appID),
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())
	assert.Empty(t, respPublicKey.GetPublicKey())

	username := randomUsername()
	pass := randomFakePassword()

	log.Println(pass)

	// проверка регистрации на ошибки
	respReg, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	realUsername := username
	realPass := pass
	realUserID := respReg.GetUserId()

	// повторная регистрация с тем же username
	respReg, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: randomFakePassword(),
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, authBusiness.ErrExists.Error())

	// проверка аутентификации на ошибки
	// проверка на неправильный appID
	respLogin, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: realUsername,
		Password: realPass,
		AppId:    int32(appID),
	})
	require.Error(t, err)
	assert.Empty(t, respLogin.GetAccessToken())
	assert.Empty(t, respLogin.GetRefreshToken())
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	// создаем приложение и производим аутентификацию
	appID, err = st.CreateApp("test")
	require.NoError(t, err)
	assert.NotEmpty(t, appID)

	// аутентификация с неправильным username
	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: randomUsername(),
		Password: randomFakePassword(),
		AppId:    int32(appID),
	})
	require.Error(t, err)
	assert.Empty(t, respLogin.GetAccessToken())
	assert.Empty(t, respLogin.GetRefreshToken())
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	// аутентификация с неправильным password
	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: realUsername,
		Password: randomFakePassword(),
		AppId:    int32(appID),
	})
	require.Error(t, err)
	assert.Empty(t, respLogin.GetAccessToken())
	assert.Empty(t, respLogin.GetRefreshToken())
	assert.ErrorContains(t, err, authBusiness.ErrInvalidCredentials.Error())

	// аутентификация
	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: realUsername,
		Password: realPass,
		AppId:    int32(appID),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetAccessToken())
	assert.NotEmpty(t, respLogin.GetRefreshToken())

	// создаем рандомный id юзера для проверки isAdmin
	fakeUserID := randomID(1, 100_000)
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &sso.IsAdminRequest{
		UserId: int64(fakeUserID),
	})
	require.Error(t, err)
	assert.Equal(t, respIsAdmin.GetIsAdmin(), false)
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	// проверяем на реального пользователя но без прав администратора
	respIsAdmin, err = st.AuthClient.IsAdmin(ctx, &sso.IsAdminRequest{
		UserId: realUserID,
	})
	require.Error(t, err)
	assert.Equal(t, respIsAdmin.GetIsAdmin(), false)
	assert.ErrorContains(t, err, authBusiness.ErrPermission.Error())

	// создаем рандомный appID для проверки getPublicKey
	fakeAppID := randomID(1, 100_000)
	respPublicKey, err = st.AuthClient.GetPublicKey(ctx, &sso.PublicKeyRequest{
		AppId: int32(fakeAppID),
	})
	require.Error(t, err)
	assert.Empty(t, respPublicKey.GetPublicKey())
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	// совершаем аутентификацию для проверки Refresh
	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: realUsername,
		Password: realPass,
		AppId:    int32(appID),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetAccessToken())
	assert.NotEmpty(t, respLogin.GetRefreshToken())

	// удаляем все приложения из бд для проверки Refresh с несуществующим appID
	err = st.CleanupAppsData()
	require.NoError(t, err)

	respRefresh, err := st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: respLogin.GetRefreshToken(),
	})
	require.Error(t, err)
	assert.Empty(t, respRefresh.GetAccessToken())
	assert.Empty(t, respRefresh.GetRefreshToken())
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	// создаем приложение, аутентифицируемся и удаляем пользователя для
	// проверки ошибки несуществующего пользователя
	appID, err = st.CreateApp("test")
	require.NoError(t, err)

	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: realUsername,
		Password: realPass,
		AppId:    int32(appID),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetAccessToken())
	assert.NotEmpty(t, respLogin.GetRefreshToken())

	err = st.CleanupUserData()
	require.NoError(t, err)

	respRefresh, err = st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: respLogin.GetRefreshToken(),
	})
	require.Error(t, err)
	assert.Empty(t, respRefresh.GetAccessToken())
	assert.Empty(t, respRefresh.GetRefreshToken())
	assert.ErrorContains(t, err, authBusiness.ErrNotFound.Error())

	username = randomUsername()
	pass = randomFakePassword()

	// проходим повторную регистрацию и аутентификацию
	respReg, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err = st.AuthClient.Login(ctx, &sso.LoginRequest{
		Username: username,
		Password: pass,
		AppId:    int32(appID),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetAccessToken())
	assert.NotEmpty(t, respLogin.GetRefreshToken())

	// проверка на получение новой пары токенов, используя отозванный refresh token
	respLogout, err := st.AuthClient.Logout(ctx, &sso.LogoutRequest{
		RefreshToken: respLogin.GetRefreshToken(),
	})
	require.NoError(t, err)
	assert.Equal(t, respLogout.GetSuccess(), true)

	respRefresh, err = st.AuthClient.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: respLogin.GetRefreshToken(),
	})
	require.Error(t, err)
	assert.ErrorContains(t, err, authBusiness.ErrTokenRevoked.Error())
	assert.Empty(t, respRefresh.GetAccessToken())
	assert.Empty(t, respRefresh.GetRefreshToken())
}
