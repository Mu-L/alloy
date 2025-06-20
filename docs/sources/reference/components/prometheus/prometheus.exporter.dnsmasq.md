---
canonical: https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.exporter.dnsmasq/
aliases:
  - ../prometheus.exporter.dnsmasq/ # /docs/alloy/latest/reference/components/prometheus.exporter.dnsmasq/
description: Learn about prometheus.exporter.dnsmasq
labels:
  stage: general-availability
  products:
    - oss
title: prometheus.exporter.dnsmasq
---

# `prometheus.exporter.dnsmasq`

The `prometheus.exporter.dnsmasq` component embeds the [`dnsmasq_exporter`](https://github.com/google/dnsmasq_exporter) for collecting statistics from a dnsmasq server.

## Usage

```alloy
prometheus.exporter.dnsmasq "<LABEL>" {
}
```

## Arguments

You can use the following arguments with `prometheus.exporter.dnsmasq`:

| Name            | Type     | Description                                          | Default                          | Required |
| --------------- | -------- | ---------------------------------------------------- | -------------------------------- | -------- |
| `address`       | `string` | The address of the dnsmasq server.                   | `"localhost:53"`                 | no       |
| `expose_leases` | `bool`   | Expose dnsmasq leases as metrics (high cardinality). | `false`                          | no       |
| `leases_file`   | `string` | The path to the dnsmasq leases file.                 | `"/var/lib/misc/dnsmasq.leases"` | no       |

## Blocks

The `prometheus.exporter.dnsmasq` component doesn't support any blocks. You can configure this component with arguments.

## Exported fields

{{< docs/shared lookup="reference/components/exporter-component-exports.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Component health

`prometheus.exporter.dnsmasq` is only reported as unhealthy if given an invalid configuration.
In those cases, exported fields retain their last healthy values.

## Debug information

`prometheus.exporter.dnsmasq` doesn't expose any component-specific
debug information.

## Debug metrics

`prometheus.exporter.dnsmasq` doesn't expose any component-specific
debug metrics.

## Example

The following example uses a [`prometheus.scrape` component][scrape] to collect metrics from `prometheus.exporter.dnsmasq`:

```alloy
prometheus.exporter.dnsmasq "example" {
  address = "localhost:53"
}

// Configure a prometheus.scrape component to collect github metrics.
prometheus.scrape "demo" {
  targets    = prometheus.exporter.dnsmasq.example.targets
  forward_to = [prometheus.remote_write.demo.receiver]
}

prometheus.remote_write "demo" {
  endpoint {
    url = "<PROMETHEUS_REMOTE_WRITE_URL>"

    basic_auth {
      username = "<USERNAME>"
      password = "<PASSWORD>"
    }
  }
}
```

Replace the following:

- _`<PROMETHEUS_REMOTE_WRITE_URL>`_: The URL of the Prometheus `remote_write` compatible server to send metrics to.
- _`<USERNAME>`_: The username to use for authentication to the `remote_write` API.
- _`<PASSWORD>`_: The password to use for authentication to the `remote_write` API.

[scrape]: ../prometheus.scrape/

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`prometheus.exporter.dnsmasq` has exports that can be consumed by the following components:

- Components that consume [Targets](../../../compatibility/#targets-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
