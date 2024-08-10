# Prophet Security Take Home

This project manages Tor exit-nodes from multiple sources. There's a control plane to manage 
the sources that we need to monitor. If an endpoint is found in multiple sources it is not removed from the 
aggregated list until it has been removed from all managed sources. For example, if source A and source B 
both had 192.168.0.1 on their list, if source A no longer included it on its list but source B did then it 
would be in our aggregated list still. As soon as source B no longer included it then it would be removed from
our list. 

You can fetch data from the service by asking for all aggregated nodes. 

There's also a control plane for creating allow lists that can filter the resulting list based on 
CIDR blocks that we don't care to track. 

## Data Sources

| Name | Url | Recommended Period |
|:-----|:----|:----| 
| Udger | https://raw.githubusercontent.com/udger/test-data/master/CSV_data_example/tor_exit_node.csv | Any |
| dan.me.uk | https://www.dan.me.uk/torlist/?exit | 00:30:00 |


## Running

### Database

Start up the database by running docker compose. This will fetch and start a postgresql container that our service will connect to. 
The database created is called `prophet-th` (this is also the name of the user and password)

```bash
docker compose -f ./deployments/docker-compose.yaml up -d
```
You'll also need to run the migrations. I chose to use `go-migrate` to manage migrations. 

```bash
# If you need need the CLI utility run the following
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 


# To execute the migrations run 
cd server/database
migrate -database postgresql://prophet-th:prophet-th@localhost:5432/prophet-th\?sslmode=disable -path migrations up
```

### Ingestion Service

There are two main components to this. The first is an async engine that will monitor endpoints on a configured period. 

```bash
# From the root of the directory run
go run ./server/ingest/cmd/server
```

### Api Server

Then start up the api server (in another terminal session)

```bash
# From the root of the directory run
go run ./server/web/cmd/server
```

For a list of all operations that the api server can do, reference the [OpenAPI specficiation](./server/api/openapi.yaml).

For a list of example requests look at the [prophet.http](./prophet.http) file.

### Optional Mock Source Server

For development a mock server that serves static CSV files can be ran in a third terminal. It returns the static files found 
in [the static directory](./server/mock/static). To test how the system reacts to endpoints getting added or removed from 
sources you can make changes to the CSV files found in that directory and the ingestion service should pick up the changes. 

```bash
# From the root of the directory run
go run ./server/mock/cmd/server
```


## Other tools

I used a few code generators to create some of the go code found in this repository. 

- [sqlc](https://github.com/sqlc-dev/sqlc) to generate models from the migration schema and defined queries found in the [database module](./server/database)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) to generate models and handler stubs from the OpenAPI spec found in the [api module](./server/api)

