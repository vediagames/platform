package graphql

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vediagames/zeroerror"

	gatewaygraphql "github.com/vediagames/platform/gateway/graphql"
	"github.com/vediagames/platform/webproxy/graphql/generated"
)

type Resolver struct {
	gatewayResolver *gatewaygraphql.Resolver
}

type Config struct {
	GatewayResolver *gatewaygraphql.Resolver
}

func (c Config) Validate() error {
	var err zeroerror.Error

	if c.GatewayResolver == nil {
		err.Add(fmt.Errorf("empty gateway resolver"))
	}

	return err.Err()
}

func NewResolver(cfg Config) Resolver {
	if err := cfg.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return Resolver{
		gatewayResolver: cfg.GatewayResolver,
	}
}

func NewSchema(r *Resolver) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: r,
	})
}
