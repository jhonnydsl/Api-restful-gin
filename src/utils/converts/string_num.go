package converts

import (
	"fmt"
	"strconv"

	"github.com/jhonnydsl/api-restful-gin/src/utils"
)

func StringToInt(value string) (int, error) {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, utils.BadRequestError(fmt.Sprintf("Erro ao converter string para inteiro %v", value))
	}
	return valueInt, nil
}