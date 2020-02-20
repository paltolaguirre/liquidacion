package Ganancias

type CalculoRemuneracionNoHabitual struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoHabitual) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES")
}

func (cg *CalculoRemuneracionNoHabitual) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETRIBUCIONES_NO_HABITUALES", 2, cg)
}

func (cg *CalculoRemuneracionNoHabitual) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionNoHabitual) getNombre() string {
	return "Remuneración No Habitual (+)"
}