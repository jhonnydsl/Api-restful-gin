package enum

import "github.com/jhonnydsl/api-restful-gin/src/utils"

type sortOrderEnum struct{}

var SortOrder sortOrderEnum

func (sortOrderEnum) AscendingInt() int {
	return 1
}

func (sortOrderEnum) DescendingInt() int {
	return -1
}

func (sortOrderEnum) AscendingStr() string {
	return "ascending"
}

func (sortOrderEnum) DescendingStr() string {
	return "descending"
}

// Converte de string para int (enum)
func (sortOrderEnum) ConvertSortOrderEnumToInt(value string) (int, error) {
	switch value {
	case "ascending":
		return SortOrder.AscendingInt(), nil
	case "descending":
		return SortOrder.DescendingInt(), nil
	default:
		return 404, utils.BadRequestError("Valor inválido para SortOrderEnum")
	}
}

// Converte de int (enum) para string
func (sortOrderEnum) ConvertSortOrderEnumToString(value int) (string, error) {
	switch value {
	case 1:
		return SortOrder.AscendingStr(), nil
	case -1:
		return SortOrder.DescendingStr(), nil
	default:
		return "unknown", utils.BadRequestError("Valor inválido para SortOrderEnum")
	}
}