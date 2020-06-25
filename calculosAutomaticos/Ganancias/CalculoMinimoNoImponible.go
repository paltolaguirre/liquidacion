package Ganancias

type CalculoMinimoNoImponible struct {
	CalculoGanancias
}

func (cg *CalculoMinimoNoImponible) getResultInternal() float64 {
	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", "valorfijomni")
	if cg.trabajoEnFechaPatagonica() {
		valorfijoMNI = 1.22 * valorfijoMNI
	}
	mesperiodoliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
	return (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
}

func (cg *CalculoMinimoNoImponible) getResult() float64 {
	return cg.getResultOnDemandTemplate("MINIMO_NO_IMPONIBLE", 42, cg)
}

func (cg *CalculoMinimoNoImponible) getTope() *float64 {
	return nil
}

func (cg *CalculoMinimoNoImponible) getNombre() string {
	return "Mínimo no Imponible"
}

func (cg *CalculoMinimoNoImponible) getEsMostrable() bool {
	return true
}
