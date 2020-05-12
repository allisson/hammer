package grpc

import "github.com/allisson/hammer"

func createdAtFilters(createdAtGt, createdAtGte, createdAtLt, createdAtLte string) []hammer.FindFilter {
	findFilters := []hammer.FindFilter{}

	if createdAtGt != "" {
		createdAtGTFilter := hammer.FindFilter{
			FieldName: "created_at",
			Operator:  "gt",
			Value:     createdAtGt,
		}
		findFilters = append(findFilters, createdAtGTFilter)
	}
	if createdAtGte != "" {
		createdAtGTEFilter := hammer.FindFilter{
			FieldName: "created_at",
			Operator:  "gte",
			Value:     createdAtGte,
		}
		findFilters = append(findFilters, createdAtGTEFilter)
	}
	if createdAtLt != "" {
		createdAtLTFilter := hammer.FindFilter{
			FieldName: "created_at",
			Operator:  "lt",
			Value:     createdAtLt,
		}
		findFilters = append(findFilters, createdAtLTFilter)
	}
	if createdAtLte != "" {
		createdAtLTEFilter := hammer.FindFilter{
			FieldName: "created_at",
			Operator:  "lte",
			Value:     createdAtLte,
		}
		findFilters = append(findFilters, createdAtLTEFilter)
	}

	return findFilters
}
