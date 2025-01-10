package casdoor

import (
	"errors"

	"github.com/kataras/iris/v12"
)

type MustAuthenticated struct {
	Options CasdoorOptions
}

func NewMustAuthenticated() *MustAuthenticated {
	return &MustAuthenticated{}
}

var (
	// ErrTokenMissing is the error value that it's returned when
	// a token is not found based on the token extractor.
	ErrTokenMissing = errors.New("required authorization token not found")
)

// Serve the middleware's action
func (m *MustAuthenticated) Serve(ctx iris.Context) {

	// If we get here, the required token is missing
	userId := _getUserId(ctx)
	if len(userId) <= 0 {
		OnError(ctx, ErrTokenMissing)
		return
	}

	// If everything ok then call next.
	ctx.Next()
}

func _getUserId(ctx iris.Context) string {
	userId, ok := ctx.Value("userId").(string)
	if ok {
		return userId
	}
	return ""
}
