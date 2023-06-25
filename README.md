# Playing with Temporal

Bits and pieces stolen from the [hello world tutorial](https://learn.temporal.io/getting_started/go/hello_world_in_go/) and the [otel tutorial](https://opentelemetry.io/docs/instrumentation/go/getting-started/).

## Prometheus Notes

Based on this [tutorial](https://docs.temporal.io/kb/prometheus-grafana-setup).

- [Prometheus](http://localhost:9090/)
- [Grafana](http://localhost:8085/)
    - [Dashboards](https://github.com/temporalio/dashboards)
    - Query examplex:
        ```promql
        sum(rate(temporal_workflow_completed_total{namespace=~\"$Namespace\"}[5m]))
        sum by (namespace) (rate(temporal_request_failure[5m]))
        sum by (operation) (rate(service_requests{service_name="frontend"}[2m]))
        histogram_quantile(0.95, sum(rate(temporal_activity_schedule_to_start_latency_bucket{}[5m])) by (namespace, activity_type, le))
        ```
- [Temporal](http://localhost:8080/)
- [Jaeger](http://localhost:16686/)