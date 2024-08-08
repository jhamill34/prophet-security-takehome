package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

type SourceResource struct {
	queries *database.Queries
}

func NewSourceResource(queries *database.Queries) *SourceResource {
	return &SourceResource{
		queries,
	}
}

func (s *SourceResource) Path() string {
	return "/sources"
}

func (s *SourceResource) Handler() http.Handler {
	r := chi.NewRouter()
	r.Post("/", s.CreateSource())
	r.Get("/", s.ListSources())
	r.Get("/{id}", s.ListSourcesNodes())
	r.Post("/{id}/start", s.StartSource())
	r.Post("/{id}/stop", s.StopSource())
	return r
}

func (s *SourceResource) ListSources() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		after := ParseIntDefault(req.URL.Query().Get("after"), -1)
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		result, err := s.queries.ListAllSources(req.Context(), database.ListAllSourcesParams{
			ID:    after,
			Limit: limit,
		})
		if err != nil {
			panic(err)
		}

		resultBytes, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(200)
		resp.Write(resultBytes)
	}
}

func (s *SourceResource) ListSourcesNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		sourceId := AssertInt(chi.URLParam(req, "id"))
		after := ParseIp(req.URL.Query().Get("after"))
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		result, err := s.queries.ListSourcesNodes(req.Context(), database.ListSourcesNodesParams{
			IpAddr: after,
			Limit:  limit,
			ID:     sourceId,
		})

		if err != nil {
			panic(err)
		}

		resultBytes, err := json.Marshal(result)
		resp.WriteHeader(200)
		resp.Write(resultBytes)
	}
}

type CreateSourceInput struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Period string `json:"period"`
}

func (s *SourceResource) CreateSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var input CreateSourceInput
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			panic(err)
		}

		var period pgtype.Interval
		period.Scan(input.Period)

		source, err := s.queries.CreateSource(req.Context(), database.CreateSourceParams{
			Name:   input.Name,
			Url:    input.Url,
			Period: period,
		})

		if err != nil {
			panic(err)
		}

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(201)
		_, err = resp.Write(sourceBytes)
		if err != nil {
			panic(err)
		}
	}
}

func (s *SourceResource) StartSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id := AssertInt(chi.URLParam(req, "id"))

		_, err := s.queries.StartSource(req.Context(), id)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}
}

func (s *SourceResource) StopSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id := AssertInt(chi.URLParam(req, "id"))

		_, err := s.queries.StopSource(req.Context(), id)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}
}
