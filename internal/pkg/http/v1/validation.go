package v1

import "fmt"

var validFieldsToOrderBy = map[string]struct{}{
	"":         {},
	"volume":   {},
	"results":  {},
	"updated":  {},
	"cpc":      {},
	"url":      {},
	"position": {},
	"keyword":  {},
}

func validateOrderByField(orderBy string) error {
	if _, ok := validFieldsToOrderBy[orderBy]; !ok {
		return fmt.Errorf("positions can't be ordered by '%s' field", orderBy)
	}

	return nil
}
