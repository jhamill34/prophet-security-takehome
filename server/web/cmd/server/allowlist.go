package main

import (
	"encoding/json"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
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
	r.Get("/", a.ListAllLists())
	r.Post("/", a.CreateAllowList())
	r.Delete("/{id}", a.DeleteAllowList())

	r.Get("/{id}/entry", a.ListAllowList())
	r.Post("/{id}/entry", a.AddToList())
	r.Delete("/{id}/entry/{entryId}", a.RemoveFromList())

	return r
}

func (a *AllowListResource) ListAllLists() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		after := ParseIntDefault(req.URL.Query().Get("after"), -1)
		limit := ParseIntDefault(req.URL.Query().Get("limit"), 10)

		dbResult, err := a.queries.ListAllLists(req.Context(), database.ListAllListsParams{
			ID:    after,
			Limit: limit,
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

		Json(req, resp, result, 200)
	}
}

func (a *AllowListResource) ListAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		dbResult, err := a.queries.ListEntriesForAllowList(req.Context(), id)
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

func (a *AllowListResource) DeleteAllowList() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		id, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "ID should be an integer", 400, err)
			return
		}

		err = a.queries.DeleteAllowList(req.Context(), int32(id))
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
		listId, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "ID should be an integer", 400, err)
			return
		}

		var input ListEntryInput
		err = json.NewDecoder(req.Body).Decode(&input)
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
			ListID: listId,
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
		listId, err := AssertInt(chi.URLParam(req, "id"))
		if err != nil {
			Err(req, resp, "Expected ID to be an integer", 400, err)
			return
		}

		entryId, err := AssertInt(chi.URLParam(req, "entryId"))
		if err != nil {
			Err(req, resp, "Expected entryID to be an integer", 400, err)
			return
		}

		err = a.queries.RemoveFromAllowlist(req.Context(), database.RemoveFromAllowlistParams{
			ListID: listId,
			ID:     entryId,
		})
		if err != nil {
			InternalServerError(req, resp, err)
		}

		resp.WriteHeader(204)
	}
}
