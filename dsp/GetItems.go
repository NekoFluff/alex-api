package dsp

import "addi/models"

func GetItems() []*models.DSPItem {
	m := make([]*models.DSPItem, 0, len(itemMap))
	for _, val := range itemMap {
		m = append(m, val)
	}
	return m
}
