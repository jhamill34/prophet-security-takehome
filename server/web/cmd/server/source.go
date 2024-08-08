package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	r.With(PaginatedContext).Get("/", s.ListSources())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(s.SourceContext)
		r.With(PaginatedContext).Get("/", s.ListSourcesNodes())
		r.Post("/start", s.StartSource())
		r.Post("/stop", s.StopSource())
	})
	return r
}

func (s *SourceResource) ListSources() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		pagination := req.Context().Value("pagination").(PaginatedInput)
		after := ParseIntDefault(pagination.Cursor, -1)

		dbResult, err := s.queries.ListAllSources(req.Context(), database.ListAllSourcesParams{
			ID:    after,
			Limit: pagination.Limit,
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

		paginated := MakePaginated(result, int(pagination.Limit), func(entry SourceEntry) string {
			return fmt.Sprintf("%d", entry.ID)
		})

		Json(req, resp, paginated, 200)
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

func (s *SourceResource) SourceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		source, err := s.queries.GetSource(req.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			Err(req, resp, "Source Not Found", 404, err)
			return
		}

		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		ctx := context.WithValue(req.Context(), "source", source)
		next.ServeHTTP(resp, req.WithContext(ctx))
	})
}

func (s *SourceResource) ListSourcesNodes() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		source := req.Context().Value("source").(database.Source)

		pagination := req.Context().Value("pagination").(PaginatedInput)
		after := ParseIp(pagination.Cursor)

		dbResult, err := s.queries.ListSourcesNodes(req.Context(), database.ListSourcesNodesParams{
			IpAddr: after,
			Limit:  pagination.Limit,
			ID:     source.ID,
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

		paginated := MakePaginated(result, int(pagination.Limit), func(entry NodeEntry) string {
			return entry.IpAddr
		})

		Json(req, resp, paginated, 200)
	}
}

func (s *SourceResource) StartSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		source := req.Context().Value("source").(database.Source)

		_, err := s.queries.StartSource(req.Context(), source.ID)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		resp.WriteHeader(204)
	}
}

func (s *SourceResource) StopSource() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		source := req.Context().Value("source").(database.Source)

		_, err := s.queries.StopSource(req.Context(), source.ID)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		resp.WriteHeader(204)
	}
}
