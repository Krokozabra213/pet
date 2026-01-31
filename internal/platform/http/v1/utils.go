package httpv1

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getUserDataCtx(c *gin.Context) (UserShortData, error) {
	value, ex := c.Get(userDataCtx)
	if !ex {
		return UserShortData{}, errors.New("userdata is missing from ctx")
	}

	userData, ok := value.(UserShortData)
	if !ok {
		return UserShortData{}, errors.New("failed to convert value from ctx to userShortData")
	}

	return userData, nil
}

func getUserIDCtx(c *gin.Context) (int, error) {
	value, ex := c.Get(userIDCtx)
	if !ex {
		return 0, errors.New("user_id is missing from ctx")
	}

	id, ok := value.(int)
	if !ok {
		return 0, errors.New("failed to convert value from ctx to int")
	}

	return id, nil
}

type UserShortData struct {
	UserID   int
	Username string
}
