package officesController

import (
	"errors"
	"fmt"
	officesError "go_template_v3/pkg/services/offices/error"
	officesModel "go_template_v3/pkg/services/offices/model"
	officesScript "go_template_v3/pkg/services/offices/script"
	"net/http"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"
)

func GetBranches(c fiber.Ctx) error {
	// Get institution code from query parameter
	instiCode := c.Query("insti_code")

	if instiCode == "" {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400, "Institution code is required.", nil, http.StatusBadRequest,
		)
	}

	// Validate institution code format (basic check)
	if len(instiCode) > 20 {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400, "Institution code is too long.", nil, http.StatusBadRequest,
		)
	}

	// Call service
	branches, err := officesScript.GetBranches(instiCode)
	if err != nil {
		// Handle specific error cases
		switch {
		case errors.Is(err, officesError.ErrInstiCodeRequired):
			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_400, "Institution code is required.", nil, http.StatusBadRequest,
			)

		case errors.Is(err, officesError.ErrNoBranchesFound):
			// Return empty array instead of error if no branches found
			response := officesModel.GetBranchesResponse{
				Branches: []officesModel.Branch{},
			}
			return v1.JSONResponseWithData(
				c, respcode.SUC_CODE_200, "No branches found.", response, http.StatusOK,
			)

		case errors.Is(err, officesError.ErrInvalidInput):
			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_400, "Invalid institution code.", nil, http.StatusBadRequest,
			)

		default:
			// Log the actual error for debugging
			fmt.Printf("Get branches error: %v\n", err)

			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_500, "Failed to fetch branches.", err, http.StatusInternalServerError,
			)
		}
	}

	// Prepare response
	response := officesModel.GetBranchesResponse{
		Branches: branches,
	}

	// Success
	return v1.JSONResponseWithData(
		c, respcode.SUC_CODE_200, "Branches fetched successfully!", response, http.StatusOK,
	)
}

func GetUnits(c fiber.Ctx) error {
	// Get branch code from query parameter
	branchCode := c.Query("branch_code")

	if branchCode == "" {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400, "Branch code is required.", nil, http.StatusBadRequest,
		)
	}

	// Validate branch code format (basic check)
	if len(branchCode) > 20 {
		return v1.JSONResponseWithError(
			c, respcode.ERR_CODE_400, "Branch code is too long.", nil, http.StatusBadRequest,
		)
	}

	// Call service
	units, err := officesScript.GetUnits(branchCode)
	if err != nil {
		// Handle specific error cases
		switch {
		case errors.Is(err, officesError.ErrBranchCodeRequired):
			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_400, "Branch code is required.", nil, http.StatusBadRequest,
			)

		case errors.Is(err, officesError.ErrInvalidInput):
			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_400, "Invalid branch code.", nil, http.StatusBadRequest,
			)

		case errors.Is(err, officesError.ErrNoUnitsFound):
			// Return empty array instead of error
			response := officesModel.GetUnitsResponse{
				Units: []officesModel.Unit{},
			}
			return v1.JSONResponseWithData(
				c, respcode.SUC_CODE_200, "No units found.", response, http.StatusOK,
			)

		default:
			// Log the actual error for debugging
			fmt.Printf("Get units error: %v\n", err)

			return v1.JSONResponseWithError(
				c, respcode.ERR_CODE_500, "Failed to fetch units.", err, http.StatusInternalServerError,
			)
		}
	}

	// Prepare response
	response := officesModel.GetUnitsResponse{
		Units: units,
	}

	// Success
	return v1.JSONResponseWithData(
		c, respcode.SUC_CODE_200, "Units fetched successfully!", response, http.StatusOK,
	)
}
