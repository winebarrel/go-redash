package redash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/winebarrel/redash-go/internal/util"
)

type QueryPage struct {
	Count    int     `json:"count"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Results  []Query `json:"results"`
}

type Query struct {
	APIKey            string          `json:"api_key"`
	CanEdit           bool            `json:"can_edit"`
	CreatedAt         time.Time       `json:"created_at"`
	DataSourceID      int             `json:"data_source_id"`
	Description       string          `json:"description"`
	ID                int             `json:"id"`
	IsArchived        bool            `json:"is_archived"`
	IsDraft           bool            `json:"is_draft"`
	IsFavorite        bool            `json:"is_favorite"`
	IsSafe            bool            `json:"is_safe"`
	LastModifiedBy    *User           `json:"last_modified_by"`
	LastModifiedByID  int             `json:"last_modified_by_id"`
	LatestQueryDataID int             `json:"latest_query_data_id"`
	Name              string          `json:"name"`
	Options           QueryOptions    `json:"options"`
	Query             string          `json:"query"`
	QueryHash         string          `json:"query_hash"`
	RetrievedAt       time.Time       `json:"retrieved_at"`
	Runtime           float64         `json:"runtime"`
	Schedule          *QueueSchedule  `json:"schedule"`
	Tags              []string        `json:"tags"`
	UpdatedAt         time.Time       `json:"updated_at"`
	User              User            `json:"user"`
	Version           int             `json:"version"`
	Visualizations    []Visualization `json:"visualizations"`
}

type QueryOptions struct {
	Parameters []map[string]any `json:"parameters"`
}

type QueueSchedule struct {
	DayOfWeek string `json:"day_of_week"`
	Interval  int    `json:"interval"`
	Time      string `json:"time"`
	Until     string `json:"until"`
}

type ListQueriesInput struct {
	OnlyFavorites bool
	Page          int
	PageSize      int
}

// https://github.com/getredash/redash-toolbelt/blob/f6d2c40881fcacb411665c75f3afbe570533539d/redash_toolbelt/client.py#L17
func (client *Client) ListQueries(ctx context.Context, input *ListQueriesInput) (*QueryPage, error) {
	params := map[string]string{}

	if input != nil {
		params["page"] = strconv.Itoa(input.Page)
		params["page_size"] = strconv.Itoa(input.PageSize)
		params["only_favorites"] = strconv.FormatBool(input.OnlyFavorites)
	}

	res, err := client.Get(ctx, "api/queries", params)

	if err != nil {
		return nil, err
	}

	page := &QueryPage{}

	if err := util.UnmarshalBody(res, &page); err != nil {
		return nil, err
	}

	return page, nil
}

// https://github.com/getredash/redash-toolbelt/blob/f6d2c40881fcacb411665c75f3afbe570533539d/redash_toolbelt/client.py#L36
func (client *Client) GetQuery(ctx context.Context, id int) (*Query, error) {
	res, err := client.Get(ctx, fmt.Sprintf("api/queries/%d", id), nil)

	if err != nil {
		return nil, err
	}

	query := &Query{}

	if err := util.UnmarshalBody(res, &query); err != nil {
		return nil, err
	}

	return query, nil
}

// https://github.com/getredash/redash-toolbelt/blob/f6d2c40881fcacb411665c75f3afbe570533539d/redash_toolbelt/client.py#L24
func (client *Client) CreateFavoriteQuery(ctx context.Context, id int) error {
	_, err := client.Post(ctx, fmt.Sprintf("api/queries/%d/favorite", id), nil)

	if err != nil {
		return err
	}

	return nil
}

// https://github.com/getredash/redash/blob/5cf13afafe4a13c8db9da645d9667bc26fd118c5/redash/handlers/queries.py#L334-L339
type CreateQueryInput struct {
	DataSourceID int                       `json:"data_source_id"`
	Description  string                    `json:"description,omitempty"`
	Name         string                    `json:"name"`
	Options      *CreateQueryInputOptions  `json:"options,omitempty"`
	Query        string                    `json:"query"`
	Schedule     *CreateQueryInputSchedule `json:"schedule,omitempty"`
}

type CreateQueryInputOptions struct {
	Parameters []map[string]any `json:"parameters"`
}

type CreateQueryInputSchedule struct {
	Interval  int     `json:"interval"`
	Time      *string `json:"time"`
	Until     *string `json:"until"`
	DayOfWeek *string `json:"day_of_week"`
}

// https://github.com/getredash/redash/blob/5cf13afafe4a13c8db9da645d9667bc26fd118c5/redash/handlers/queries.py#L207-L212
func (client *Client) CreateQuery(ctx context.Context, input *CreateQueryInput) (*Query, error) {
	res, err := client.Post(ctx, "api/queries", input)

	if err != nil {
		return nil, err
	}

	query := &Query{}

	if err := util.UnmarshalBody(res, &query); err != nil {
		return nil, err
	}

	return query, nil
}

// https://github.com/getredash/redash-toolbelt/blob/f6d2c40881fcacb411665c75f3afbe570533539d/redash_toolbelt/client.py#L130
func (client *Client) ForkQuery(ctx context.Context, id int) (*Query, error) {
	res, err := client.Post(ctx, fmt.Sprintf("api/queries/%d/fork", id), nil)

	if err != nil {
		return nil, err
	}

	query := &Query{}

	if err := util.UnmarshalBody(res, &query); err != nil {
		return nil, err
	}

	return query, nil
}

// https://github.com/getredash/redash-toolbelt/blob/f6d2c40881fcacb411665c75f3afbe570533539d/redash_toolbelt/client.py#L147
type UpdateQueryInput struct {
	DataSourceID int                       `json:"data_source_id,omitempty"`
	Description  string                    `json:"description,omitempty"`
	Name         string                    `json:"name,omitempty"`
	Options      *UpdateQueryInputOptions  `json:"options,omitempty"`
	Query        string                    `json:"query,omitempty"`
	Schedule     *UpdateQueryInputSchedule `json:"schedule,omitempty"`
}

type UpdateQueryInputOptions struct {
	Parameters []map[string]any `json:"parameters"`
}

type UpdateQueryInputSchedule struct {
	Interval  int     `json:"interval"`
	Time      *string `json:"time"`
	Until     *string `json:"until"`
	DayOfWeek *string `json:"day_of_week"`
}

// https://github.com/getredash/redash/blob/5cf13afafe4a13c8db9da645d9667bc26fd118c5/redash/handlers/queries.py#L207-L212
func (client *Client) UpdateQuery(ctx context.Context, id int, input *UpdateQueryInput) (*Query, error) {
	res, err := client.Post(ctx, fmt.Sprintf("api/queries/%d", id), input)

	if err != nil {
		return nil, err
	}

	query := &Query{}

	if err := util.UnmarshalBody(res, &query); err != nil {
		return nil, err
	}

	return query, nil
}

func (client *Client) ArchiveQuery(ctx context.Context, id int) error {
	_, err := client.Delete(ctx, fmt.Sprintf("api/queries/%d", id))

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetQueryResultsJSON(ctx context.Context, id int) ([]byte, error) {
	return client.GetQueryResults(ctx, id, "json")
}

func (client *Client) GetQueryResultsCSV(ctx context.Context, id int) ([]byte, error) {
	return client.GetQueryResults(ctx, id, "csv")
}

func (client *Client) GetQueryResults(ctx context.Context, id int, ext string) ([]byte, error) {
	res, err := client.Get(ctx, fmt.Sprintf("api/queries/%d/results.%s", id, ext), nil)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

func (client *Client) ExecQueryJSON(ctx context.Context, id int) ([]byte, string, error) {
	res, err := client.Post(ctx, fmt.Sprintf("api/queries/%d/results", id), map[string]string{"filetype": "json"})

	if err != nil {
		return nil, "", err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, "", err
	}

	if bytes.HasPrefix(body, []byte(`{"job": {"status": 1, "error": "", "id": "`)) {
		var job JobResponse
		err := json.Unmarshal(body, &job)

		if err == nil {
			return body, job.Job.ID, nil
		}
	}

	return body, "", nil
}
