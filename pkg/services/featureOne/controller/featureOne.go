package ctrFeatureOne

import (
	"go_template_v3/pkg/config"
	mdlFeatureOne "go_template_v3/pkg/services/featureOne/model"
	scpFeatureOne "go_template_v3/pkg/services/featureOne/script"
	"net/http"

	v1 "github.com/FDSAP-Git-Org/hephaestus/helper/v1"
	"github.com/FDSAP-Git-Org/hephaestus/respcode"
	"github.com/gofiber/fiber/v3"

)


func GetItems(c fiber.Ctx) error {
	var item []mdlFeatureOne.ItemBody
	var cat mdlFeatureOne.CategoryParams

	script := scpFeatureOne.AddProduct
	c.Bind().Body(&cat)

	config.DBConnList[0].Debug().Raw(script, cat.Category).Scan(&item)
	return v1.JSONResponseWithData(c, respcode.SUC_CODE_200, respcode.SUC_CODE_200_MSG, item, http.StatusOK)
}
