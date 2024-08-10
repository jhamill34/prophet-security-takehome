package routes

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

// ListAggregatedNodes implements api.StrictServerInterface.
func (s *ServerRoutes) ListAggregatedNodes(ctx context.Context, request api.ListAggregatedNodesRequestObject) (api.ListAggregatedNodesResponseObject, error) {
	resultMap := make(map[string]*api.NodeEntry, 0)
	invert := DefaultValue(request.Params.Invert, false)
	limit := DefaultValue(request.Params.Limit, 10)
	after, err := ParseIp(request.Params.After)
	if err != nil {
		return nil, err
	}

	if request.Params.AllowlistId == nil {
		dbResult, err := s.queries.ListAllNodes(ctx, database.ListAllNodesParams{
			IpAddr: after,
			Limit:  int32(limit),
		})

		if err != nil {
			// TODO: Come back here...
			return nil, err
		}

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
	} else {
		if invert {
			dbResult, err := s.queries.ListNodesWithoutAllowlist(ctx, database.ListNodesWithoutAllowlistParams{
				IpAddr: after,
				Limit:  int32(limit),
				ListID: int32(*request.Params.AllowlistId),
			})
			if err != nil {
				// TODO: ...
				return nil, err
			}

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
		} else {
			dbResult, err := s.queries.ListFilteredAllowlistNodes(ctx, database.ListFilteredAllowlistNodesParams{
				IpAddr: after,
				Limit:  int32(limit),
				ListID: int32(*request.Params.AllowlistId),
			})
			if err != nil {
				// TODO: ....
				return nil, err
			}

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
		}

	}

	result := make([]api.NodeEntry, 0, len(resultMap))
	for _, v := range resultMap {
		result = append(result, *v)
	}

	slices.SortFunc(result, func(a, b api.NodeEntry) int {
		return strings.Compare(a.IpAddr, b.IpAddr)
	})

	paginatedMetadata := MakePaginated(result, limit, func(item api.NodeEntry) string {
		return item.IpAddr
	})

	response := api.ListAggregatedNodes200JSONResponse{
		Cursor:  paginatedMetadata.Cursor,
		HasMore: paginatedMetadata.HasMore,
		Total:   paginatedMetadata.Total,
		Data:    result,
	}

	return response, nil
}
