package httpv1

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getUserDataCtx(c *gin.Context) (UserShortData, error) {
	value, ex := c.Get(userData)
	if !ex {
		return UserShortData{}, errors.New("userdata is missing from ctx")
	}

	userData, ok := value.(UserShortData)
	if !ok {
		return UserShortData{}, errors.New("failed to convert value from ctx to userShortData")
	}

	return userData, nil
}

type UserShortData struct {
	UserID   string
	Username string
}
