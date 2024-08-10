package routes

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

// CreateSource implements api.StrictServerInterface.
func (s *ServerRoutes) CreateSource(ctx context.Context, request api.CreateSourceRequestObject) (api.CreateSourceResponseObject, error) {
	var period pgtype.Interval
	err := period.Scan(request.Body.Period)
	if err != nil {
		return nil, err
	}

	dbResult, err := s.queries.CreateSource(ctx, database.CreateSourceParams{
		Name:   request.Body.Name,
		Url:    request.Body.Url,
		Period: period,
	})

	if err != nil {
		return nil, err
	}

	resultPeriod, err := dbResult.Period.Value()
	if err != nil {
		return nil, err
	}
	result := api.SourceEntry{
		Id:            int(dbResult.ID),
		Name:          dbResult.Name,
		Url:           dbResult.Url,
		Period:        resultPeriod.(string),
		LastExecution: dbResult.LastExecution.Time.Format(time.RFC3339),
		Version:       int(dbResult.Version.Int64),
		Running:       dbResult.Running.Bool,
	}

	return api.CreateSource201JSONResponse(result), nil
}

// ListSourceNodes implements api.StrictServerInterface.
func (s *ServerRoutes) ListSourceNodes(ctx context.Context, request api.ListSourceNodesRequestObject) (api.ListSourceNodesResponseObject, error) {
	source, err := s.queries.GetSource(ctx, int32(request.Id))
	if err != nil {
		return nil, err
	}

	limit := DefaultValue(request.Params.Limit, 10)
	after, err := ParseIp(request.Params.After)
	if err != nil {
		return nil, err
	}

	dbResult, err := s.queries.ListSourcesNodes(ctx, database.ListSourcesNodesParams{
		IpAddr: after,
		Limit:  int32(limit),
		ID:     source.ID,
	})
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]*api.NodeEntry)
	for _, r := range dbResult {
		sourceEntry := api.NodeSourceEntry{
			SourceId:      int(r.SourceID),
			Version:       int(r.Version.Int64),
			LastExecution: r.LastExecution.Time.Format(time.RFC3339),
		}
		if entry, ok := resultMap[r.IpAddr.String()]; ok {
			entry.Sources = append(entry.Sources, sourceEntry)
		} else {
			resultMap[r.IpAddr.String()] = &api.NodeEntry{
				IpAddr:  r.IpAddr.String(),
				Sources: []api.NodeSourceEntry{sourceEntry},
			}
		}
	}

	result := make([]api.NodeEntry, 0, len(resultMap))
	for _, v := range resultMap {
		result = append(result, *v)
	}

	slices.SortFunc(result, func(a, b api.NodeEntry) int {
		return strings.Compare(a.IpAddr, b.IpAddr)
	})

	paginatedMetadata := MakePaginated(result, limit, func(entry api.NodeEntry) string {
		return entry.IpAddr
	})

	paginated := api.PaginatedNodeEntry{
		Total:   paginatedMetadata.Total,
		HasMore: paginatedMetadata.HasMore,
		Cursor:  paginatedMetadata.Cursor,
		Data:    result,
	}

	return api.ListSourceNodes200JSONResponse(paginated), nil
}

// ListSources implements api.StrictServerInterface.
func (s *ServerRoutes) ListSources(ctx context.Context, request api.ListSourcesRequestObject) (api.ListSourcesResponseObject, error) {
	limit := DefaultValue(request.Params.Limit, 10)
	after, err := strconv.ParseInt(DefaultValue(request.Params.After, "-1"), 10, 32)
	if err != nil {
		// TODO: ....
		return nil, err
	}

	dbResult, err := s.queries.ListAllSources(ctx, database.ListAllSourcesParams{
		ID:    int32(after),
		Limit: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]api.SourceEntry, len(dbResult))
	for i, r := range dbResult {
		period, err := r.Period.Value()
		if err != nil {
			return nil, err
		}

		result[i] = api.SourceEntry{
			Id:            int(r.ID),
			Name:          r.Name,
			Url:           r.Url,
			Period:        period.(string),
			LastExecution: r.LastExecution.Time.Format(time.RFC3339),
			Version:       int(r.Version.Int64),
			Running:       r.Running.Bool,
		}
	}

	paginatedMetadata := MakePaginated(result, limit, func(entry api.SourceEntry) string {
		return fmt.Sprintf("%d", entry.Id)
	})

	paginated := api.PaginatedSourceEntry{
		Cursor:  paginatedMetadata.Cursor,
		Total:   paginatedMetadata.Total,
		HasMore: paginatedMetadata.HasMore,
		Data:    result,
	}

	return api.ListSources200JSONResponse(paginated), nil
}

// StartSource implements api.StrictServerInterface.
func (s *ServerRoutes) StartSource(ctx context.Context, request api.StartSourceRequestObject) (api.StartSourceResponseObject, error) {
	source, err := s.queries.GetSource(ctx, int32(request.Id))

	_, err = s.queries.StartSource(ctx, source.ID)
	if err != nil {
		return nil, err
	}

	return api.StartSource204Response{}, nil
}

// StopSource implements api.StrictServerInterface.
func (s *ServerRoutes) StopSource(ctx context.Context, request api.StopSourceRequestObject) (api.StopSourceResponseObject, error) {
	source, err := s.queries.GetSource(ctx, int32(request.Id))

	_, err = s.queries.StopSource(ctx, source.ID)
	if err != nil {
		return nil, err
	}

	return api.StopSource204Response{}, nil
}
