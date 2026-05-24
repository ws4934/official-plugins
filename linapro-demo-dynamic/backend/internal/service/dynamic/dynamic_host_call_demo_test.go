// This file verifies host-call demo helpers that do not require a running Wasm
// host service.

package dynamicservice

import "testing"

// fakeConfigHostService returns deterministic plugin config values for unit
// tests.
type fakeConfigHostService struct {
	strings map[string]configStringResult
	bools   map[string]configBoolResult
}

// configStringResult stores one string config read result.
type configStringResult struct {
	value string
	found bool
}

// configBoolResult stores one bool config read result.
type configBoolResult struct {
	value bool
	found bool
}

// String returns one configured fake string value.
func (s *fakeConfigHostService) String(key string) (string, bool, error) {
	result := s.strings[key]
	return result.value, result.found, nil
}

// Bool returns one configured fake bool value.
func (s *fakeConfigHostService) Bool(key string) (bool, bool, error) {
	result := s.bools[key]
	return result.value, result.found, nil
}

// fakeHostConfigHostService returns deterministic public host config values for unit
// tests.
type fakeHostConfigHostService struct {
	strings map[string]configStringResult
	bools   map[string]configBoolResult
}

// String returns one configured fake public host config string value.
func (s *fakeHostConfigHostService) String(key string) (string, bool, error) {
	result := s.strings[key]
	return result.value, result.found, nil
}

// Bool returns one configured fake public host config bool value.
func (s *fakeHostConfigHostService) Bool(key string) (bool, bool, error) {
	result := s.bools[key]
	return result.value, result.found, nil
}

// TestRunHostCallDemoConfigReadsPluginAndHostConfigValues verifies the dynamic
// demo reads plugin config and public host config values through separate host
// service clients.
func TestRunHostCallDemoConfigReadsPluginAndHostConfigValues(t *testing.T) {
	service := &serviceImpl{
		configSvc: &fakeConfigHostService{
			strings: map[string]configStringResult{
				hostCallDemoPluginGreetingKey: {
					value: "Hello from test config",
					found: true,
				},
			},
			bools: map[string]configBoolResult{
				hostCallDemoPluginFeatureKey: {
					value: true,
					found: true,
				},
			},
		},
		hostConfigSvc: &fakeHostConfigHostService{
			strings: map[string]configStringResult{
				hostCallDemoWorkspaceKey: {
					value: "/tmp/linapro",
					found: true,
				},
				hostCallDemoI18nDefaultKey: {
					value: "zh-CN",
					found: true,
				},
			},
			bools: map[string]configBoolResult{
				hostCallDemoI18nEnabledKey: {
					value: true,
					found: true,
				},
			},
		},
	}

	payload, err := service.runHostCallDemoConfig()
	if err != nil {
		t.Fatalf("expected config demo to succeed, got error: %v", err)
	}
	if !payload.Plugin.GreetingFound || payload.Plugin.Greeting != "Hello from test config" {
		t.Fatalf("unexpected plugin greeting payload: %#v", payload.Plugin)
	}
	if !payload.Plugin.FeatureEnabledFound || !payload.Plugin.FeatureEnabled {
		t.Fatalf("unexpected plugin feature payload: %#v", payload.Plugin)
	}
	if !payload.HostConfig.WorkspaceBasePathFound || payload.HostConfig.WorkspaceBasePath != "/tmp/linapro" {
		t.Fatalf("unexpected host workspace payload: %#v", payload.HostConfig)
	}
	if !payload.HostConfig.I18nDefaultFound || payload.HostConfig.I18nDefault != "zh-CN" {
		t.Fatalf("unexpected host i18n default payload: %#v", payload.HostConfig)
	}
	if !payload.HostConfig.I18nEnabledFound || !payload.HostConfig.I18nEnabled {
		t.Fatalf("unexpected host i18n enabled payload: %#v", payload.HostConfig)
	}
}
