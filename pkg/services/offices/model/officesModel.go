package officesModel

// ----------------------------
// BRANCH MODEL
// ----------------------------
type Branch struct {
	BranchCode string `json:"branch_code"`
	BranchName string `json:"branch_name"`
}

type GetBranchesResponse struct {
	Branches []Branch `json:"branches"`
}

type GetBranchesRequest struct {
	InstiCode string `json:"insti_code"`
}

// ----------------------------
// UNIT MODEL
// ----------------------------

type Unit struct {
	UnitCode string `json:"unit_code"`
	UnitName string `json:"unit_name"`
}

type GetUnitsResponse struct {
	Units []Unit `json:"units"`
}

type GetUnitsRequest struct {
	BranchCode string `json:"branch_code"`
}
