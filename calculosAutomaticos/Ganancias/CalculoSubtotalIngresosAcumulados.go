package Ganancias

import "github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"

type CalculoSubtotalIngresosAcumulados struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalIngresosAcumulados) getResultInternal() float64 {

	var importeTotal float64 = 0

	if cg.esTipoSac() {
		importeTotal = cg.getSac(true)
		//Solo lo ejecuto, pero no lo sumo
		(&CalculoSubtotalIngresos{cg.CalculoGanancias}).getResult()
	} else {
		importeTotal = importeTotal + (&CalculoSubtotalIngresos{cg.CalculoGanancias}).getResult()

		var arraySubtotalDeduccionesGenerales []float64

		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidios{cg.CalculoGanancias}).getResult())
		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesObraSocial{cg.CalculoGanancias}).getResult())
		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoCuotaSindical{cg.CalculoGanancias}).getResult())
		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos{cg.CalculoGanancias}).getResult())
		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesObraSocialOtrosEmpleos{cg.CalculoGanancias}).getResult())
		arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoCuotaSindicalOtrosEmpleos{cg.CalculoGanancias}).getResult())

		importeTotal = importeTotal - Sum(arraySubtotalDeduccionesGenerales)
	}
	var itemGananciaAnterior *structLiquidacion.Liquidacionitem
	liquidacionAnterior := cg.obtenerLiquidacionIgualAnioLegajoMesAnterior()

	if liquidacionAnterior != nil {
		itemGananciaAnterior = obtenerItemGananciaFromLiquidacion(liquidacionAnterior, cg.Db)
	}


	if itemGananciaAnterior != nil {
		importeTotal = importeTotal + (&CalculoSubtotalIngresosAcumulados{CalculoGanancias{itemGananciaAnterior, liquidacionAnterior, cg.Db, false}}).getResult()
	}

	return importeTotal
}

func (cg *CalculoSubtotalIngresosAcumulados) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_INGRESOS_ACUMULADOS", 22, cg)
}

func (cg *CalculoSubtotalIngresosAcumulados) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalIngresosAcumulados) getNombre() string {
	return "Subtotal Ingresos Acumulados"
}

func (cg *CalculoSubtotalIngresosAcumulados) getEsMostrable() bool {
	return true
}
