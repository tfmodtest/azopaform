package pkg

import (
	"fmt"
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"json-rule-finder/pkg/shared"
	"strings"

	aztypes "github.com/ms-henglu/go-azure-types/types"
)

var azureTypes = func() sets.Set {
	var r = hashset.New()
	for name, _ := range aztypes.DefaultAzureSchemaLoader().GetSchema().Resources {
		r.Add(name)
	}
	return r
}()

type OperationField string

func (o OperationField) Rego(ctx *shared.Context) (string, error) {
	return o.processedFieldName()
}

func (o OperationField) processedFieldName() (string, error) {
	s := string(o)
	if !strings.Contains(s, "/") {
		return s, nil
	}
	split := strings.Split(s, "/")
	expectedAzureType := strings.Join(split[:len(split)-1], "/")
	if !azureTypes.Contains(expectedAzureType) {
		return "", fmt.Errorf("unknown azure type: %s", expectedAzureType)
	}
	propertyPath := split[len(split)-1]
	propertyPath = strings.ReplaceAll(propertyPath, "[*]", "[_]")
	return fmt.Sprintf("r.change.after.%s", propertyPath), nil
}
