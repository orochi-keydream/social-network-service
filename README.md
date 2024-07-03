# social-network-service

* [1 Prerequisites](#1-prerequisites)
* [2 Prepare the environment](#2-prepare-the-environment)
* [3 Run the service](#3-run-the-service)
  * [3.1 Run locally](#31-run-locally)
  * [3.2 Run using Docker](#32-run-using-docker)
* [4 Explore the functionality](#4-explore-the-functionality)
* [5 Activity simulation](#5-activity-simulation)
* [6 Notifications via WebSockets](#6-notifications-via-websockets)

## 1 Prerequisites

* [Go](https://go.dev/) (`v1.22` or later)
* [Docker](https://www.docker.com/)
* [goose](https://github.com/pressly/goose) (for running migrations)

`goose` can be installed by the following command (Go language must be already installed on the machine):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

`goose` can also be installed using `brew` (for MacOS):

```bash
brew install goose
```

## 2 Prepare the environment

#### Step 1 - Up containers

Run the following command to up all infrastructure containers.

```bash
docker compose -f ./docker-compose-infra.yml up -d
```

#### Step 2 - Prepare database

Using http://localhost:7000/, determine which PostgreSQL node is master and run the following commands (replace `postgres0` if needed).

```bash
docker exec -it postgres0 bash
psql -U postgres
create database social_network_db;
exit
```

#### Step 3 - Configure Patroni cluster

Staying inside the container, run the command below to change Patroni configuration.

```bash
patronictl edit-config
```

Add the following line at the end of the configuration file.

```
synchronous_mode: on
```

Save and exit from the editor (type `:wq`, then press `Enter`). Accept the changes to apply them to Patroni cluster.

Make sure that master, sync and async replicas are available on http://localhost:7000/. Note that it may take some time before changes are applied on the page.

#### Step 4 - Run migrations

Execute the command below to apply migrations:

```bash
goose -dir ./migrations/ postgres "host=localhost port=15432 user=postgres password=123 dbname=social_network_db" up
```

#### Step 5 - Configure topics (optional)

> NOTE
> 
> This step is optional. All topics will be created automatically with default settings (the number of partitions in this case will be equal to 1 for each topic), but you can configure topics manually to configure the number of partitions and other settings. You will also be able to read messages, monitor consumer groups and perform some administrative operations.

Go to http://localhost:8082/ and press `Configure new cluster`. Specify any cluster name and confiugre Kafka servers (`kafka1:29091`, `kafka2:29092`, `kafka3:29093`). Press `Validate` button to make sure that everything is OK, then press `Submit` button to get access to the cluster.

Create `post_events` and `feed_cache_commands` topics with desired settings.

## 3 Run the service

Depending on your preferences, you can run the service using one of the following ways:

* Run locally (Go must be installed)
* Run using Docker

### 3.1 Run locally

Run the following command:

```bash
go run ./cmd/app/main.go --config ./config/local.yml
```

### 3.2 Run using Docker

Run the following command:

```bash
docker compose -f ./docker-compose-service.yml up -d
```

## 4 Explore the functionality

Open http://localhost:8080/swagger/index.html to work with the service via Swagger.

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

Check other endpoints using Swagger to explore other features of the service.

## 5 Activity simulation

This repository has an utility that simulates activity on the service. Run the following to start the simulation (Go is required):

```bash
go run ./cmd/sim/main.go 10 ./scripts/
```

* The first argument (`10`) is a number of goroutines (users) that work with the service concurrently.
* The last argument (`./scripts/`) is a directory that must contain the following files: `cities.txt`, `male-names.txt`, `male-names.txt`, `surnames.txt` and `posts.txt`. These files are used to generate user profiles and their posts. You can change the content of these files or completely replace them, but you should follow the format of the files (all strings are separated with one line-break character).

At the moment users can:

* Create posts
* Read feed
* Add friends
* Remove friends

## 6 Notifications via WebSockets

http://localhost:8081/hub endpoint is used to connect to the service using WebSocket that allows you to be notified in realtime. You can use Postman or other tools to check, how it works.

> NOTE
> 
> `/hub` endpoint requires authorization, so you need to add `Authorization` header with `Bearer {token}` value, where `{token}` is a JWT token that you are given when you call `/login` endpoint.

Currently, the service sends the following types of messages:

* A new post appeared
* An existing post updated

The short description of the process is described [here](https://github.com/orochi-keydream/social-network-service/wiki/WebSocket-scheme-description).
