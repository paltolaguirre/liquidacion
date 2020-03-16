package Ganancias

type CalculoRemuneracionBruta struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionBruta) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", false)
}

func (cg *CalculoRemuneracionBruta) getResult() float64 {
	return cg.getResultOnDemandTemplate("REMUNERACION_BRUTA", 1, cg)
}

func (cg *CalculoRemuneracionBruta) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionBruta) getNombre() string {
	return "Remuneración Bruta (+)"
}

func (cg *CalculoRemuneracionBruta) getEsMostrable() bool {
	return true
}
