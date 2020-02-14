package Ganancias

import "fmt"

type CalculoGananciaNeta struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNeta) getResultInternal() float64 {
	var arrayGananciaNeta []float64
	var gananciaNeta float64

	arrayGananciaNeta = append(arrayGananciaNeta, (&CalculoCuotaMedicoAsistencial{cg.CalculoGanancias}).getResult())
	arrayGananciaNeta = append(arrayGananciaNeta, (&CalculoDonacionFisicosNacProvMunArt20{cg.CalculoGanancias}).getResult())

	gananciaNeta = (&CalculoSubtotal{cg.CalculoGanancias}).getResult() - Sum(arrayGananciaNeta)
	fmt.Println("Calculos Automaticos - Ganancia Neta:", gananciaNeta)
	return gananciaNeta
}

func (cg *CalculoGananciaNeta) getResult() float64 {
	return cg.getResultOnDemandTemplate("Ganancia Neta", "GANANCIA_NETA", 38, cg)
}

func (cg *CalculoGananciaNeta) getTope() *float64 {
	return nil
}
