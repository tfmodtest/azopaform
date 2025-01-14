package pkg

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

func ResourceTypeParser(resourceType string) (string, error) {
	upperRt := strings.ToUpper(resourceType)
	ttt, ok := shared.ResourceTypeLookupTable.QueryResource(upperRt, "")
	if !ok || len(ttt) == 0 {
		/*
			return "", fmt.Errorf("cannot find the resource type %s in the lookup table", resourceType)
		*/
		return "", nil
	}
	var result string
	for _, v := range ttt {
		result = v[0].ResourceType
		break
	}
	// The `azurerm_app_service_plan` resource has been superseded by the `azurerm_service_plan` resource.
	if result == "azurerm_app_service_plan" {
		result = "azurerm_service_plan"
	} else if result == "azurerm_app_service_environment" {
		result = "azurerm_app_service_environment_v3"
	} else if result == "azurerm_sql_server" {
		result = "azurerm_mssql_server"
	}
	return result, nil
}
