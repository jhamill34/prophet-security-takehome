package routes

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

func (s *ServerRoutes) getAllowList(ctx context.Context, id int32) (api.AllowlistEntry, error) {
	dbResult, err := s.queries.GetAllowList(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return api.AllowlistEntry{}, fmt.Errorf("Allow list not found")
	}

	if err != nil {
		return api.AllowlistEntry{}, err
	}

	result := api.AllowlistEntry{
		Id:   int(dbResult.ID),
		Name: dbResult.Name,
	}

	return result, nil
}

// AddToAllowlist implements api.StrictServerInterface.
func (s *ServerRoutes) AddToAllowlist(ctx context.Context, request api.AddToAllowlistRequestObject) (api.AddToAllowlistResponseObject, error) {
	list, err := s.getAllowList(ctx, int32(request.Id))
	if err != nil {
		return nil, err
	}

	ipAddr, err := netip.ParsePrefix(request.Body.Cidr)
	if err != nil {
		return nil, err
	}

	dbResult, err := s.queries.AddToAllowlist(ctx, database.AddToAllowlistParams{
		Cidr:   ipAddr,
		ListID: int32(list.Id),
	})
	if err != nil {
		return nil, err
	}

	entry := api.AllowlistEntryItem{
		Id:          int(dbResult.ID),
		Cidr:        dbResult.Cidr.String(),
		AllowlistId: int(dbResult.ListID),
	}

	return api.AddToAllowlist201JSONResponse(entry), nil
}

// CreateAllowlist implements api.StrictServerInterface.
func (s *ServerRoutes) CreateAllowlist(ctx context.Context, request api.CreateAllowlistRequestObject) (api.CreateAllowlistResponseObject, error) {
	dbResult, err := s.queries.CreateAllowList(ctx, request.Body.Name)
	if err != nil {
		return nil, err
	}

	entry := api.AllowlistEntry{
		Id:   int(dbResult.ID),
		Name: dbResult.Name,
	}

	return api.CreateAllowlist201JSONResponse(entry), nil
}

// DeleteAllowList implements api.StrictServerInterface.
func (s *ServerRoutes) DeleteAllowList(ctx context.Context, request api.DeleteAllowListRequestObject) (api.DeleteAllowListResponseObject, error) {
	panic("unimplemented")
}

// ListAllAllowlists implements api.StrictServerInterface.
func (s *ServerRoutes) ListAllAllowlists(ctx context.Context, request api.ListAllAllowlistsRequestObject) (api.ListAllAllowlistsResponseObject, error) {
	after, err := strconv.ParseInt(DefaultValue(request.Params.After, "-1"), 10, 32)
	if err != nil {
		// TODO: ....
		return nil, err
	}
	limit := DefaultValue(request.Params.Limit, 10)

	dbResult, err := s.queries.ListAllLists(ctx, database.ListAllListsParams{
		ID:    int32(after),
		Limit: int32(limit),
	})

	if err != nil {
		// TODO:
		return nil, err
	}

	result := make([]api.AllowlistEntry, len(dbResult))
	for i, r := range dbResult {
		result[i] = api.AllowlistEntry{
			Id:   int(r.ID),
			Name: r.Name,
		}
	}

	paginated := MakePaginated(result, limit, func(item api.AllowlistEntry) string {
		return fmt.Sprintf("%d", item.Id)
	})

	response := api.PaginatedAllowlistEntry{
		Total:   paginated.Total,
		Cursor:  paginated.Cursor,
		HasMore: paginated.HasMore,
		Data:    result,
	}

	return api.ListAllAllowlists200JSONResponse(response), nil
}

// ListAllowlistEntries implements api.StrictServerInterface.
func (s *ServerRoutes) ListAllowlistEntries(ctx context.Context, request api.ListAllowlistEntriesRequestObject) (api.ListAllowlistEntriesResponseObject, error) {
	list, err := s.getAllowList(ctx, int32(request.Id))
	if err != nil {
		return nil, err
	}
	dbResult, err := s.queries.ListEntriesForAllowList(ctx, int32(list.Id))
	if err != nil {
		return nil, err
	}

	entries := make([]api.AllowlistEntryItem, len(dbResult))

	for i, r := range dbResult {
		entries[i] = api.AllowlistEntryItem{
			Id:          int(r.ID),
			Cidr:        r.Cidr.String(),
			AllowlistId: int(r.ListID),
		}
	}

	return api.ListAllowlistEntries200JSONResponse(entries), nil
}

// RemoveFromAllowlist implements api.StrictServerInterface.
func (s *ServerRoutes) RemoveFromAllowlist(ctx context.Context, request api.RemoveFromAllowlistRequestObject) (api.RemoveFromAllowlistResponseObject, error) {
	list, err := s.getAllowList(ctx, int32(request.Id))
	if err != nil {
		return nil, err
	}

	err = s.queries.RemoveFromAllowlist(ctx, database.RemoveFromAllowlistParams{
		ListID: int32(list.Id),
		ID:     int32(request.EntryId),
	})
	if err != nil {
		return nil, err
	}

	return api.RemoveFromAllowlist204Response{}, nil
}
