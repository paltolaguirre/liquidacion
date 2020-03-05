package Ganancias

type CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica struct {
	CalculoGanancias
}

func (cg *CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "HONORARIOS_SERV_ASISTENCIA_SANITARIA_MEDICA_Y_PARAMEDICA", "deducciondesgravacionsiradig")
	return importeTotal
}

func (cg *CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica) getResult() float64 {
	return cg.getResultOnDemandTemplate("HONORARIOS_SERV_ASISTENCIA_SANITARIA_MEDICA_Y_PARAMEDICA", 0, cg)
}

func (cg *CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica) getTope() *float64 {
	return nil
}

func (cg *CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica) getNombre() string {
	return "Honorarios Serv. Asistencia Sanitaria Medica y Paramedica"
}

func (cg *CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica) getEsMostrable() bool {
	return false
}
