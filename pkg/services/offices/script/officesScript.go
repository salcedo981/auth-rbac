package officesScript

import (
	"encoding/json"
	"fmt"
	"go_template_v3/pkg/config"
	officesError "go_template_v3/pkg/services/offices/error"
	officesModel "go_template_v3/pkg/services/offices/model"
	"strings"
)

func GetBranches(instiCode string) ([]officesModel.Branch, error) {
	db := &config.DBConnList[0]

	// Validate input
	if instiCode == "" {
		return nil, officesError.ErrInstiCodeRequired
	}

	// Call the PostgreSQL function
	query := `
		SELECT get_branches(?)::text
	`

	var jsonText string
	err := db.Raw(query, instiCode).Scan(&jsonText).Error

	if err != nil {
		// Log the error for debugging
		fmt.Printf("Get branches database error: %v\n", err)
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}

	// Parse the JSON text into slice of Branch
	var branches []officesModel.Branch
	err = json.Unmarshal([]byte(jsonText), &branches)
	if err != nil {
		return nil, fmt.Errorf("failed to parse branches response: %w", err)
	}

	// If no branches found, return empty slice (not error)
	// The PostgreSQL function already returns empty array '[]'
	if branches == nil {
		branches = []officesModel.Branch{}
	}

	return branches, nil
}

func GetUnits(branchCode string) ([]officesModel.Unit, error) {
	db := &config.DBConnList[0]

	// Validate input
	if branchCode == "" {
		return nil, officesError.ErrBranchCodeRequired
	}

	// Call the PostgreSQL function
	query := `
		SELECT get_units(?)::text
	`

	var jsonText string
	err := db.Raw(query, branchCode).Scan(&jsonText).Error

	if err != nil {
		// Check for specific database errors
		errMsg := err.Error()

		// Handle potential foreign key or constraint errors
		if strings.Contains(errMsg, "violates foreign key constraint") {
			return nil, fmt.Errorf("invalid branch code: %w", officesError.ErrInvalidInput)
		}

		// Log the error for debugging
		fmt.Printf("Get units database error: %v\n", err)
		return nil, fmt.Errorf("failed to fetch units: %w", err)
	}

	// Parse the JSON text into slice of Unit
	var units []officesModel.Unit
	err = json.Unmarshal([]byte(jsonText), &units)
	if err != nil {
		return nil, fmt.Errorf("failed to parse units response: %w", err)
	}

	// If no units found, return empty slice (not error)
	if units == nil {
		units = []officesModel.Unit{}
	}

	return units, nil
}
