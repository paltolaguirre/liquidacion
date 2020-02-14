package Ganancias

type CalculoRemuneracionBruta struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionBruta) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA")
}

func (cg *CalculoRemuneracionBruta) getResult() float64 {
	return cg.getResultOnDemandTemplate("Remuneracion Bruta", "REMUNERACION_BRUTA", 1, cg)
}

func (cg *CalculoRemuneracionBruta) getTope() *float64 {
	return nil
}