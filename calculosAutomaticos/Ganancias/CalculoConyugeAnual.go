package Ganancias

type CalculoConyugeAnual struct {
	CalculoGanancias
}

func (cg *CalculoConyugeAnual) getResultInternal() float64 {
	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", "valorfijomni")
	return cg.getfgDetalleCargoFamiliarAnual("conyugeid", "valorfijoconyuge", 100, valorfijoMNI)
}

func (cg *CalculoConyugeAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("CONYUGE_ANUAL", 0, cg)
}

func (cg *CalculoConyugeAnual) getTope() *float64 {
	return nil
}

func (cg *CalculoConyugeAnual) getNombre() string {
	return "Conyuge Anual"
}

func (cg *CalculoConyugeAnual) getEsMostrable() bool {
	return false
}
