package Ganancias

type CalculoRetribucionesNoHabitualesExentasNoAlcanzadas struct {
	CalculoGanancias
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES_EXENTAS_NO_ALCANZADAS")
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETRIBUCIONES_NO_HABITUALES_EXENTAS_NO_ALCANZADAS", 0, cg)
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadas) getTope() *float64 {
	return nil
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadas) getNombre() string {
	return "Retribuciones No Habituales Exentas/No Alcanzadas"
}

func (cg *CalculoRetribucionesNoHabitualesExentasNoAlcanzadas) getEsMostrable() bool {
	return false
}
