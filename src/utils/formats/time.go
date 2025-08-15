package formats

import (
	"fmt"
	"time"

	"github.com/jhonnydsl/api-restful-gin/src/utils"
)

func Time(t time.Time, format string) (time.Time, error) {
	timeStr := t.Format(format)
	timeFormatted, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, utils.BadRequestError(fmt.Sprintf("Erro ao formatar data: %v", err))
	}
	return timeFormatted, nil
}