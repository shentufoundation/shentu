package types

const (
	QueryProgramByID = "program_id"
)

// QueryPaginationParams provides basic pagination parameters
// for queries in bounty module.
type QueryPaginationParams struct {
	Page  int
	Limit int
}

// NewQueryPaginationParams creates new instance of the
// QueryPaginationParams.
func NewQueryPaginationParams(page, limit int) QueryPaginationParams {
	return QueryPaginationParams{
		Page:  page,
		Limit: limit,
	}
}
