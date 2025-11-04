package appmiddleware

import (
	"context"
	"go-restapi/pkg/common"
	"go-restapi/pkg/common/token"
)

func GetUserFromContext(ctx context.Context) (*token.Payload, error) {
	payload, ok := ctx.Value(UserContextKey).(*token.Payload)
	if !ok {
		return nil, common.ErrUnauthorized
	}
	return payload, nil
}
