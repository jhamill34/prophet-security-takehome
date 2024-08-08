package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

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

		dbResult, err := s.queries.ListAllSources(req.Context(), database.ListAllSourcesParams{
			ID:    after,
			Limit: limit,
		})
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		result := make([]SourceEntry, len(dbResult))
		for i, r := range dbResult {
			period, err := r.Period.Value()
			if err != nil {
				InternalServerError(req, resp, err)
				return
			}

			result[i] = SourceEntry{
				ID:            r.ID,
				Name:          r.Name,
				Url:           r.Url,
				Period:        period.(string),
				LastExecution: r.LastExecution.Time.Format(time.RFC3339),
				Version:       r.Version.Int64,
				Running:       r.Running.Bool,
			}
		}

		Json(req, resp, result, 200)
	}
}

func (s *SourceResource) ListSourcesNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		sourceId, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		after := ParseIp(req.URL.Query().Get("after"))
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		dbResult, err := s.queries.ListSourcesNodes(req.Context(), database.ListSourcesNodesParams{
			IpAddr: after,
			Limit:  limit,
			ID:     sourceId,
		})
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		resultMap := make(map[string]*NodeEntry)
		for _, r := range dbResult {
			sourceEntry := NodeSourceEntry{
				SourceId:      r.SourceID,
				Version:       r.Version.Int64,
				LastExecution: r.LastExecution.Time.Format(time.RFC3339),
			}
			if entry, ok := resultMap[r.IpAddr.String()]; ok {
				entry.Sources = append(entry.Sources, sourceEntry)
			} else {
				resultMap[r.IpAddr.String()] = &NodeEntry{
					IpAddr:  r.IpAddr.String(),
					Sources: []NodeSourceEntry{sourceEntry},
				}
			}
		}

		result := make([]NodeEntry, 0, len(resultMap))
		for _, v := range resultMap {
			result = append(result, *v)
		}

		slices.SortFunc(result, func(a, b NodeEntry) int {
			return strings.Compare(a.IpAddr, b.IpAddr)
		})

		Json(req, resp, result, 200)
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
			Err(req, resp, "Unable to parse JSON input", 400, err)
			return
		}

		var period pgtype.Interval
		err = period.Scan(input.Period)
		if err != nil {
			Err(req, resp, "Invalid format for period should be in the form '[y years] [m mons] [d days] HH:MM:SS' where [] are optional", 400, err)
			return
		}

		dbResult, err := s.queries.CreateSource(req.Context(), database.CreateSourceParams{
			Name:   input.Name,
			Url:    input.Url,
			Period: period,
		})

		if err != nil {
			panic(err)
		}

		resultPeriod, err := dbResult.Period.Value()
		if err != nil {
			panic(err)
		}
		result := SourceEntry{
			ID:            dbResult.ID,
			Name:          dbResult.Name,
			Url:           dbResult.Url,
			Period:        resultPeriod.(string),
			LastExecution: dbResult.LastExecution.Time.Format(time.RFC3339),
			Version:       dbResult.Version.Int64,
			Running:       dbResult.Running.Bool,
		}

		Json(req, resp, result, 201)
	}
}

func (s *SourceResource) StartSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		_, err = s.queries.StartSource(req.Context(), id)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		resp.WriteHeader(204)
	}
}

func (s *SourceResource) StopSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		_, err = s.queries.StopSource(req.Context(), id)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		resp.WriteHeader(204)
	}
}
