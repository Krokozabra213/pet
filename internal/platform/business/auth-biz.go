package business

type AuthService struct {
	// TODO: ADD CLIENT FOR SSO SERVICE
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	return "", "", nil
}
