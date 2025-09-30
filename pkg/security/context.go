package security

import (
	"context"
	"fmt"
)

func GetIDFromContext(ctx context.Context) (int, error) {
	something := ctx.Value("userID")
	userID, ok := something.(int)
	if !ok {
		return 0, fmt.Errorf("failed get id from context")
	}
	return userID, nil
}
