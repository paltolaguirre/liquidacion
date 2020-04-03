package Ganancias

type CalculoBaseImponible struct {
	CalculoGanancias
}

func (cg *CalculoBaseImponible) getResultInternal() float64 {
	var arrayBaseImponible []float64
	var importeTotal float64

	arrayBaseImponible = append(arrayBaseImponible, (&CalculoGananciaNetaAcumSujetaAImp{cg.CalculoGanancias}).getResult())
	arrayBaseImponible = append(arrayBaseImponible, (&CalculoDeduccionesPersonales{cg.CalculoGanancias}).getResult())

	importeTotal = Sum(arrayBaseImponible)
	return importeTotal
}

func (cg *CalculoBaseImponible) getResult() float64 {
	return cg.getResultOnDemandTemplate("BASE_IMPONIBLE", 48, cg)
}

func (cg *CalculoBaseImponible) getTope() *float64 {
	return nil
}

func (cg *CalculoBaseImponible) getNombre() string {
	return "Base Imponible"
}

func (cg *CalculoBaseImponible) getEsMostrable() bool {
	return true
}
