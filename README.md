# UES Projection Service

This service consumes the primary application log for the Universal
Event Store and creates a common projection of the current state. This
state can then be queried by other services via a simple REST API.

## Running with Docker Compose

Copy the sample .env files from `env/sample` to `env/`.

```sh
cp env/sample/*.env env/
```

Then, add your IAM credentials to each of the .env files. These
credentials should have permission to write and read objects from
the `ues-events` bucket.

Finally, you should be able to run this service locally with a simple
`docker-compose up`.
