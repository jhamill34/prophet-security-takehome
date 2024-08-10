package auth

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func AuthValidator(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	return nil
}
