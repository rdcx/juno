package auth

import (
	"context"
	"juno/pkg/api/user"
)

// Define a key to avoid context key collisions
type contextKey string

const userContextKey = contextKey("authUser")

// Store the user in the context
func WithUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// Retrieve the user from the context
func UserFromContext(ctx context.Context) (*user.User, bool) {
	user, ok := ctx.Value(userContextKey).(*user.User)
	return user, ok
}

func MustUserFromContext(ctx context.Context) *user.User {
	user, ok := UserFromContext(ctx)
	if !ok {
		panic("user not found in context")
	}
	return user
}
