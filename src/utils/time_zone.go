package utils

import "time"

func TimeNowBrazil() time.Time {
	brazilLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		InternalServerError("Erro ao carregar fuso horario de Brasilia: " + err.Error())
	}
	return time.Now().In(brazilLocation)
}

func TimeBrazil(valueTime time.Time) time.Time {
	brazilLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		InternalServerError("Erro ao carregar fuso horario de Brasilia: " + err.Error())
	}
	return valueTime.In(brazilLocation)
}