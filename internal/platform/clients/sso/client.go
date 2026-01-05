package ssoclient

import (
	"context"
	"fmt"
	"time"

	"github.com/Krokozabra213/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api     sso.AuthClient
	appID   int32
	timeout time.Duration
}

func NewClient(timeout time.Duration, address string, appID int32) (*Client, error) {
	const op = "clients.sso.newclient"

	retryOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.NewClient(address, retryOpts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api:     sso.NewAuthClient(cc),
		appID:   appID,
		timeout: timeout,
	}, nil
}

func (c *Client) GetPublicKey(ctx context.Context) (string, error) {
	const op = "clients.sso.getpublickey"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.GetPublicKey(ctx, &sso.PublicKeyRequest{
		AppId: c.appID,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.PublicKey, nil
}

func (c *Client) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "clients.sso.isadmin"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.IsAdmin(ctx, &sso.IsAdminRequest{
		UserId: userID,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.IsAdmin, nil
}

func (c *Client) Register(ctx context.Context, username string, password string) (int64, error) {
	const op = "clients.sso.register"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.Register(ctx, &sso.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resp.UserId, nil
}

func (c *Client) Login(ctx context.Context, username string, password string) (string, string, error) {
	const op = "clients.sso.login"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.Login(ctx, &sso.LoginRequest{
		AppId:    c.appID,
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.AccessToken, resp.RefreshToken, nil
}

func (c *Client) Logout(ctx context.Context, refreshToken string) (bool, error) {
	const op = "clients.sso.logout"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.Logout(ctx, &sso.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.Success, nil
}

func (c *Client) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "clients.sso.refresh"

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.api.Refresh(ctx, &sso.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.AccessToken, resp.RefreshToken, nil
}
