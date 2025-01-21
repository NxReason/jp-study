package viewmodels

import "jp.study/m/v2/models"

func RadicalList(m map[int]models.Radical) []models.Radical {
	radicals := make([]models.Radical, len(m))

	index := 0
	for _, v := range m {
		radicals[index] = v
		index++
	}

	return radicals
}