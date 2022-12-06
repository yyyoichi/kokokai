package ctx

import (
	"context"
	"kokokai/server/auth"
)

type UserKey string

const (
	userClaims UserKey = "claims"
)

func NewUserContext(ctx context.Context, mc *auth.MyClaims) context.Context {
	return context.WithValue(ctx, userClaims, mc)
}
func FromUserContext(ctx context.Context) (*auth.MyClaims, bool) {
	u, ok := ctx.Value(userClaims).(*auth.MyClaims)
	return u, ok
}
