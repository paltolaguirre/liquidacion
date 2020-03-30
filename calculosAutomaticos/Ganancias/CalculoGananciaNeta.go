package Ganancias

type CalculoGananciaNeta struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNeta) getResultInternal() float64 {
	var arrayGananciaNeta []float64
	var gananciaNeta float64

	arrayGananciaNeta = append(arrayGananciaNeta, (&CalculoCuotaMedicoAsistencial{cg.CalculoGanancias}).getResult())
	arrayGananciaNeta = append(arrayGananciaNeta, (&CalculoDonacionFisicosNacProvMunArt20{cg.CalculoGanancias}).getResult())

	gananciaNeta = (&CalculoSubtotal{cg.CalculoGanancias}).getResult() - Sum(arrayGananciaNeta)
	return gananciaNeta
}

func (cg *CalculoGananciaNeta) getResult() float64 {
	return cg.getResultOnDemandTemplate("GANANCIA_NETA", 39, cg)
}

func (cg *CalculoGananciaNeta) getTope() *float64 {
	return nil
}

func (cg *CalculoGananciaNeta) getNombre() string {
	return "Ganancia Neta"
}

func (cg *CalculoGananciaNeta) getEsMostrable() bool {
	return true
}
