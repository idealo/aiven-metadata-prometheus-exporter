# Aiven Metadata Prometheus Exporter

## Purpose

Provide monitoring and observability on metadata information of [Aiven](https://aiven.io/), especially account, project and service
  information.

## Available Metrics

| Metric Name                           | Description |
|---------------------------------------|---|
| aiven_account_team_count_total        | The number of teams per account|
| aiven_account_team_member_count_total | The number of members per team for an account|
| aiven_project_count_total             | The number of projects registered in the account|
| aiven_service_count_total             | The number of services per project|
| aiven_project_estimated_billing_total | The estimated billing per project|
| aiven_project_vpc_count_total         | The number of VPCs per project|
| aiven_project_vpc_peering_count_total | The number of VPC peering connections per project|
| aiven_service_node_count_total        | Node count per service|
| aiven_service_node_state_info         | Node state per service|
| aiven_service_serviceuser_count_total | Service user count per service|
| aiven_service_topic_count_total       | Topic count per service|
| aiven_service_booked_plan_info        | The booked plan for a service|

## Usage

This prometheus exporter leverages the [Aiven Go SDK](https://github.com/aiven/aiven-go-client)

The following arguments are available:

    Usage of bin/aiven-metadata-prometheus-exporter:
      -debug
            Enable debug logging
      -listen-address string
            Address to listen on for telemetry (default ":2112")
      -scrape-interval string
            Aiven API scrape interval (default "5m")
      -telemetry-path string
            Path under which to expose metrics (default "/metrics")


## Contributing

*Contributions are highly welcome!*

* Feel free to contribute enhancements or bug fixes.
  * Fork this repo, apply your changes and create a PR pointing to this repo and the develop branch
* If you have any ideas or suggestions, please open an issue and describe your idea or feature request
  * If you have any more generic requests, feel free to open a [discussion](https://github.com/idealo/aiven-metadata-prometheus-exporter/discussions) 

## License

This project is licensed under the MIT License - see the LICENSE file for details
