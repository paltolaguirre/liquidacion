package Ganancias

type CalculoRemuneracionNoHabitual struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoHabitual) getResultInternal() float64 {
	total := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES")
	return total + cg.obtenerConceptosProrrateoMesesAnteriores()
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

func (cg *CalculoRemuneracionNoHabitual) getEsMostrable() bool {
	return true
}
