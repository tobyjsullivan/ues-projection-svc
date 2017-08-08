# UES Projection Service

This service consumes the primary application log for the Universal
Event Store and creates a common projection of the current state. This
state can then be queried by other services via a simple REST API.

## Running with Docker Compose

You should be able to run this service locally with a simple
`docker-compose up`.
