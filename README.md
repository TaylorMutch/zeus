# Zeus

Add-ons for Grafana LGTM stack that help facilitate operating an OSS observability stack

## Why zeus?

Zeus was born out of a need to enable OSS deployments of Grafana's LGTM (Loki / Grafana / Tempo / Mimir) stack to be functional as a consistent set of backends in a multi-tenant environment.

Some of the problems solved by Zeus include:

* Providing a simple authentication layer which tightly binds to LGTM's multi-tenant model
* Providing plumbing for combining multiple tenant alerts into a single alertmanager backend, while still disambiguating between tenants in alerts
* Automating datasource creation in Grafana for one or more backends, and configuring them to enable correlations

Additionally, Zeus provides components for:

* Storing additional objects associated with a tenant, such as HTTP/TCP/UDP probes that can be consumed by [blackbox exporter](https://github.com/prometheus/blackbox_exporter)

## Contributing

See the [contribution guide](./CONTRIBUTING.md) for further details!
