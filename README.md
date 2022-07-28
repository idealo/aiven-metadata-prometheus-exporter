# Aiven Prometheus Exporter

## Work in progress!

### Purpose

* Provide monitoring and observability for the Cloud Shuttle team to see how Aiven behaves and performs

* The following metrics are collected and exposed via a prometheus endpoint on :2112/metrics:
  * Projects Count
  * Service (cluster) count per project
  * Node count per service
  * Node state per service
  * Service user count per service

### Todo:
* Metrics
  * [x] No. VPC Peering Cxt per Project
  * [x] No. VPCs per Project
  * [x] No. Topics per Service 
  * [X] Booked Plan per Service
  * [x] Estimated billing per project

* Test coverage 
