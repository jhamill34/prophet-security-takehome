package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

type AllowListResource struct {
	queries *database.Queries
}

func NewAllowListResource(queries *database.Queries) *AllowListResource {
	return &AllowListResource{
		queries,
	}
}

func (a *AllowListResource) Path() string {
	return "/allowlist"
}

func (a *AllowListResource) Handler() http.Handler {
	r := chi.NewRouter()
	r.With(PaginatedContext).Get("/", a.ListAllLists())
	r.Post("/", a.CreateAllowList())

	r.Route("/{id}", func(r chi.Router) {
		r.Use(a.AllowListContext)

		r.Delete("/", a.DeleteAllowList())
		r.Get("/entry", a.ListAllowList())
		r.Post("/entry", a.AddToList())
		r.Delete("/entry/{entryId}", a.RemoveFromList())
	})

	return r
}

func (a *AllowListResource) ListAllLists() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		pagination := req.Context().Value("pagination").(PaginatedInput)
		after := ParseIntDefault(pagination.Cursor, -1)

		dbResult, err := a.queries.ListAllLists(req.Context(), database.ListAllListsParams{
			ID:    after,
			Limit: pagination.Limit,
		})

		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		result := make([]AllowlistEntry, len(dbResult))
		for i, r := range dbResult {
			result[i] = AllowlistEntry{
				ID:   r.ID,
				Name: r.Name,
			}
		}

		paginated := MakePaginated(result, int(pagination.Limit), func(item AllowlistEntry) string {
			return fmt.Sprintf("%d", item.ID)
		})

		Json(req, resp, paginated, 200)
	}
}

type CreateAllowListInput struct {
	Name string `json:"name"`
}

func (a *AllowListResource) CreateAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var input CreateAllowListInput
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			Err(req, resp, "Unable to parse JSON input", 400, err)
			return
		}

		dbResult, err := a.queries.CreateAllowList(req.Context(), input.Name)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		entry := AllowlistEntry{
			ID:   dbResult.ID,
			Name: dbResult.Name,
		}

		Json(req, resp, entry, 201)
	}
}

func (a *AllowListResource) AllowListContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		list, err := a.queries.GetAllowList(req.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			Err(req, resp, "Allow List Not Found", 404, err)
			return
		}

		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		ctx := context.WithValue(req.Context(), "allowlist", list)

		next.ServeHTTP(resp, req.WithContext(ctx))
	})
}

func (a *AllowListResource) ListAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		list := req.Context().Value("allowlist").(database.Allowlist)

		dbResult, err := a.queries.ListEntriesForAllowList(req.Context(), list.ID)
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		entries := make([]AllowlistEntryItem, len(dbResult))

		for i, r := range dbResult {
			entries[i] = AllowlistEntryItem{
				ID:     r.ID,
				Cidr:   r.Cidr.String(),
				ListID: r.ListID,
			}
		}

		Json(req, resp, entries, 200)
	}
}

func (a *AllowListResource) DeleteAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		list := req.Context().Value("allowlist").(database.Allowlist)

		err := a.queries.DeleteAllowList(req.Context(), list.ID)
		if err != nil {
			InternalServerError(req, resp, err)
		}

		resp.WriteHeader(204)
	}
}

type ListEntryInput struct {
	Cidr string `json:"cidr"`
}

func (a *AllowListResource) AddToList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		list := req.Context().Value("allowlist").(database.Allowlist)

		var input ListEntryInput
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			Err(req, resp, "Unable to parse JSON input", 400, err)
			return
		}

		ipAddr, err := netip.ParsePrefix(input.Cidr)
		if err != nil {
			Err(req, resp, "Expected Cidr be of the form <ip>/<bits>", 400, err)
			return
		}

		dbResult, err := a.queries.AddToAllowlist(req.Context(), database.AddToAllowlistParams{
			Cidr:   ipAddr,
			ListID: list.ID,
		})
		if err != nil {
			InternalServerError(req, resp, err)
			return
		}

		entry := AllowlistEntryItem{
			ID:     dbResult.ID,
			Cidr:   dbResult.Cidr.String(),
			ListID: dbResult.ListID,
		}

		Json(req, resp, entry, 201)
	}
}

func (a *AllowListResource) RemoveFromList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		list := req.Context().Value("allowlist").(database.Allowlist)

		entryId, err := AssertInt(chi.URLParam(req, "entryId"))
		if err != nil {
			Err(req, resp, "Expected entryID to be an integer", 400, err)
			return
		}

		err = a.queries.RemoveFromAllowlist(req.Context(), database.RemoveFromAllowlistParams{
			ListID: list.ID,
			ID:     entryId,
		})
		if err != nil {
			InternalServerError(req, resp, err)
		}

		resp.WriteHeader(204)
	}
}
