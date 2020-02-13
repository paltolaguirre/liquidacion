package Ganancias

type CalculoRemuneracionNoHabitual struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoHabitual) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES")
}

func (cg *CalculoRemuneracionNoHabitual) getResult() float64 {
	return cg.getResultOnDemandTemplate("Remuneracion No Habitual", "RETRIBUCIONES_NO_HABITUALES", 2, cg)
}
