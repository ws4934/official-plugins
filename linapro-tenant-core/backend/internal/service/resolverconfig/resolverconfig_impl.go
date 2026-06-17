// resolverconfig_impl.go implements resolver configuration validation,
// persistence, ordering, and cache-friendly lookup for the linapro-tenant-core plugin.
// It keeps strategy validation aligned with resolver behavior so invalid
// tenant-detection rules are rejected before runtime use.

package resolverconfig

import (
	"context"

	"lina-plugin-linapro-tenant-core/backend/internal/service/resolver"
	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// Get returns the code-owned resolver policy.
func (s *serviceImpl) Get(ctx context.Context) (*Config, error) {
	return defaultConfig(), nil
}

// ToResolverConfig projects the code-owned resolver policy into the resolver
// package config. It honors the policy chain (which includes the domain
// resolver), reserved subdomains, root domain, and ambiguity behavior. The root
// domain stays empty by default, which keeps subdomain resolution disabled until
// an operator configures one; custom-domain resolution does not depend on it.
func ToResolverConfig(config *Config) resolver.Config {
	if config == nil {
		config = defaultConfig()
	}
	return resolver.Config{
		Chain:              cloneStrings(config.Chain),
		ReservedSubdomains: cloneStrings(config.ReservedSubdomains),
		RootDomain:         config.RootDomain,
		OnAmbiguous:        config.OnAmbiguous,
	}
}

// defaultConfig returns the built-in resolver configuration.
func defaultConfig() *Config {
	return &Config{
		Chain:              shared.DefaultResolverChain(),
		ReservedSubdomains: shared.DefaultReservedSubdomains(),
		RootDomain:         shared.DefaultRootDomain,
		OnAmbiguous:        shared.OnAmbiguousPrompt,
		Version:            1,
	}
}

// cloneStrings returns a detached copy of string slices stored in the policy.
func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
