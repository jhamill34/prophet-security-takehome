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
	r.Post("/{sourceName}/start", s.StartSource())
	r.Post("/{sourceName}/stop", s.StopSource())
	return r
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
		sourceName := chi.URLParam(req, "sourceName")

		_, err := s.queries.StartSource(req.Context(), sourceName)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}
}

func (s *SourceResource) StopSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		sourceName := chi.URLParam(req, "sourceName")

		_, err := s.queries.StopSource(req.Context(), sourceName)
		if err != nil {
			panic(err)
		}

		resp.WriteHeader(204)
	}
}
