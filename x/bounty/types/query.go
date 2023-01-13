package types

// query endpoints supported by the bounty Querier
const (
	QueryHosts    = "hosts"
	QueryHost     = "host"
	QueryPrograms = "programs"
	QueryProgram  = "program"
	QueryFindings = "findings"
	QueryFinding  = "finding"
)

// QueryProgramsParams Params for query 'custom/bounty/programs'
type QueryProgramsParams struct {
	Page  int
	Limit int
	// TODO add filter
}

// NewQueryProgramsParams creates a new instance of NewQueryProgramsParams
func NewQueryProgramsParams(page, limit int) QueryProgramsParams {
	return QueryProgramsParams{
		Page:  page,
		Limit: limit,
	}
}

// QueryProgramParams Params for query 'custom/bounty/program'
type QueryProgramParams struct {
	ProgramID uint64
}

// NewQueryProgramIDParams creates a new instance of ProgramIDParams
func NewQueryProgramIDParams(programID uint64) QueryProgramParams {
	return QueryProgramParams{
		ProgramID: programID,
	}
}
