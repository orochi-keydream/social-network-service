# social-network-service

* [Prerequisites](#prerequisites)
* [Get started](#get-started)

## Prerequisites

* [Go](https://go.dev/) (`v1.22` or later)
* [Docker](https://www.docker.com/)
* [goose](https://github.com/pressly/goose) (for migrations)
* make (optional)

`goose` can be installed by the following command (Go language must be already installed on the machine):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

`goose` can also be installed using `brew` (for MacOS):

```
brew install goose
```

## Get started

Prepare a database and run the service using one of the following ways.

Using `make`:

1. Run `make compose-up`.
2. Run `make migrate-up`.
3. Run `make run`.

Or executing commands directly:

1. Run `docker compose -f docker-compose.yml up -d`.
2. Run `goose up`.
3. Run `go run ./cmd/main.go`.

Choose a way to interact with the service:

* Open http://localhost:8080/swagger/index.html to work with the service via Swagger (recommended).
* Open [Postman Collection](https://www.postman.com/orochi-keydream/workspace/public/collection/27125449-0ec721d4-3b0c-4a41-b3b9-ea79cb0ba811?action=share&creator=27125449).

Register a new user profile using `POST /user/register` endpoint with the folowing body:

```json
{
  "biography": "Software developer",
  "birthdate": "1990-01-01",
  "city": "New York",
  "first_name": "John",
  "gender": "Male",
  "password": "123456",
  "second_name": "Doe"
}
```

This handler returns a new `user_id` that uses UUID v4 format:

```json
{
  "user_id": "bbeb7da8-6d75-4419-9d94-91ec52bc506c"
}
```

Provide the given `user_id` and password to `POST /login` endpoint to get JWT token:

```json
{
  "password": "123456",
  "user_id": "bbeb7da8-6d75-4419-9d94-91ec52bc506c"
}
```

Use the given `token` to make requests that require authorization. The token must be placed to `Authorization` header and follow after `Bearer` prefix.

Try to create a new post using `POST /post/create` endpoint using the following body (note that this endpoint requires authorization):

```json
{
  "text": "Hello, World!"
}
```

Make sure that post has been created calling `GET /post/feed` without any parameters. A response may look like below:

```json
{
  "posts": [
    {
      "postId": "e9e80db7-b75c-4e73-8002-dd171026b634",
      "authorUserId": "bbeb7da8-6d75-4419-9d94-91ec52bc506c",
      "text": "Hello, World!"
    }
  ]
}
```

Check other endpoints using Swagger or Postman Collection to explore other features of the service.
