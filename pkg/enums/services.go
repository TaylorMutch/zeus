package enums

import "fmt"

var (
	AlertmanagerLabelProxyServiceName = serviceName("alertmanager-label-proxy")
	AuthServiceName                   = serviceName("auth")
)

func serviceName(name string) string {
	return fmt.Sprintf("zeus-%s", name)
}
