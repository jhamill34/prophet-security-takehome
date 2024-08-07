package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

//
// One service polls for existing ingestion configurations
// and schedules the ingestion based on the state we read
//
// Another service polls for active jobs
//

func main() {
	db := NewDatabase(context.TODO(), "host=localhost port=5432 user=prophet-th password=prophet-th dbname=prophet-th sslmode=disable")
	queries := database.New(db)

	httpClient := NewHttpClient()
	ingester := NewIngester(queries, httpClient)
	ingester.Run(context.TODO())
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}

type Ingester struct {
	queries    *database.Queries
	httpClient *http.Client
}

func NewIngester(queries *database.Queries, httpClient *http.Client) *Ingester {
	return &Ingester{
		queries,
		httpClient,
	}
}

func (i *Ingester) Run(ctx context.Context) {
	for {
		sources, err := i.queries.ListEligableSources(ctx)
		if err != nil {
			panic(err)
		}

		for _, s := range sources {
			childCtx := context.WithValue(ctx, "job_name", s.Name)
			i.doIngestion(childCtx, s)
		}

		i.idle()
	}
}

func (i *Ingester) idle() {
	slog.Info("Sleeping for 10 seconds")
	time.Sleep(10 * time.Second)
}

func (i *Ingester) doIngestion(ctx context.Context, source database.Source) {
	slog.InfoContext(ctx, "Bumping source version")
	source, err := i.queries.PrepareExecution(ctx, source.Name)
	if err != nil {
		panic(err)
	}

	slog.InfoContext(ctx, "Fetching canonical data")
	req, err := http.NewRequestWithContext(ctx, "GET", source.Url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := i.httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(resp.Body)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	var sourceId pgtype.Int4
	sourceId.Scan(int64(source.ID))

	insertNodes := make([]database.BatchInsertNodesParams, len(rows))
	for i, row := range rows {
		var newVersion pgtype.Int8
		newVersion.Scan(source.Version.Int64 + 1)

		insertNodes[i] = database.BatchInsertNodesParams{
			IpAddr:   row[0],
			SourceID: source.ID,
			Version:  newVersion,
		}
	}

	results := i.queries.BatchInsertNodes(ctx, insertNodes)
	results.Exec(func(i int, err error) {
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
		}

		slog.InfoContext(ctx, fmt.Sprintf("Inserted %d nodes", i))
	})
}

func NewDatabase(ctx context.Context, connection string) *pgx.Conn {
	db, err := pgx.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	return db
}
