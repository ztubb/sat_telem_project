# Telemetry Processing Stack

A test project implementing a minimally viable stack for ingesting and serving telemetry data.

## Design

This project consists of four primary services:
1. The mock data publisher generates telemetry data at a fixed rate of one message/sec and publishes it on request via a UDP connection. It is primarily sourced from the provided python script.
2. The Go data ingestor initiates a UDP connection and receives data from the mock data publisher. Received data is parsed and stored in a single PostgreSQL table.
3. The Go API service is colocated with the ingestor service and executed with the appropriate command parameter. It handles REST requests from clients, queries the database for the requested data, and returns it as JSON.
4. A PostgreSQL database is used to store ingested telemetry data.

Services are containerized and deployed/configured with docker compose.

## Setup

1. Clone this repository.
2. Copy the `.env.example` into a `.env` file and make your local changes
3. To build and run the stack, execute the following command
   - ```docker compose up --build```

## Scalability Improvements
- To handle increased data rates, creating a "go-between" service to pull all data from the publisher down into a single global work queue could suffice. The Go ingestor clients would then need to be updated to process data from that global queue. This will enable us to scale up the number of ingestor replicas to a sufficient level where all incoming data is processed in near-real time. 
- Observability and monitoring improvements should be implemented to stop ingesting data if any of the required services experience outages (e.g. if the database has gone down, we no longer want to process incoming data). Much of this could be implementing by moving the container orchestration into Kubernetes and leveraging `livenessProbe` configurations.
- Leveraging Kubernetes resources (services) will also enable us to increase replica counts for the API to enable high-availability and increased request traffic.

## Additional Improvements
- Additional configurations to change telemetry producer data rates
- API documentation and error-handling
- API caching to reduce round-trip Postgres queries
- Websocket support in API layer to enable server to client comms (alerts)
- Database partitioning/sharding to improve query performance, e.g. partitioning by day so that all telemetry for a given day is quickly accessible.