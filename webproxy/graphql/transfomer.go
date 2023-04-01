package graphql

import (
	"github.com/vediagames/platform/gateway/graphql/model"
)

func sortingMethodToPointer(m model.SortingMethod) *model.SortingMethod {
	c := m
	return &c
}
