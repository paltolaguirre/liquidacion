package Ganancias

type CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importeconceptosexentos")
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETRIBUCIONES_NO_HABITUALES_EXENTAS_NO_ALCANZADAS_OTROS_EMPLEOS", 0, cg)
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos) getNombre() string {
	return "Retribuciones No Habituales Exentas / No Alcanzadas Otros Empleos"
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadasOtrosEmpleos) getEsMostrable() bool {
	return false
}
