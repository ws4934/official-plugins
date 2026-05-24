// This file binds the sample plugin service to the guest-side pluginbridge
// host-service clients. The framework guest SDK provides real clients for
// wasip1 builds and unsupported stubs for ordinary Go test builds.

package dynamicservice

import "lina-core/pkg/pluginbridge"

// newRuntimeHostService returns the guest-side runtime host client.
func newRuntimeHostService() runtimeHostService {
	return pluginbridge.Runtime()
}

// newStorageHostService returns the guest-side storage host client.
func newStorageHostService() storageHostService {
	return pluginbridge.Storage()
}

// newNetworkHostService returns the guest-side outbound HTTP host client.
func newNetworkHostService() networkHostService {
	return pluginbridge.HTTP()
}

// newCronHostService returns the guest-side cron registration host client.
func newCronHostService() cronHostService {
	return pluginbridge.Cron()
}

// newConfigHostService returns the guest-side plugin config host client.
func newConfigHostService() configHostService {
	return pluginbridge.Config()
}

// newHostConfigHostService returns the guest-side public host config client.
func newHostConfigHostService() hostConfigHostService {
	return pluginbridge.HostConfig()
}
