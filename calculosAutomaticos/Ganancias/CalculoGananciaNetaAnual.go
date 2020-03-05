package Ganancias

import "fmt"

type CalculoGananciaNetaAnual struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNetaAnual) getResultInternal() float64 {
	var gananciaNetaAnual float64

	gananciaNetaAnual = (&CalculoSubtotalAnual{cg.CalculoGanancias}).getResult() - (&CalculoCuotaMedicoAsistencial{cg.CalculoGanancias}).getResult() - (&CalculoDonacionFisicosNacProvMunArt20{cg.CalculoGanancias}).getResult()
	fmt.Println("Calculos Automaticos - Ganancia Neta Anual:", gananciaNetaAnual)
	return gananciaNetaAnual
}

func (cg *CalculoGananciaNetaAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("GANANCIA_NETA_ANUAL", 0, cg)
}

func (cg *CalculoGananciaNetaAnual) getTope() *float64 {
	return nil
}

func (cg *CalculoGananciaNetaAnual) getNombre() string {
	return "Ganancia Neta Anual"
}

func (cg *CalculoGananciaNetaAnual) getEsMostrable() bool {
	return false
}
