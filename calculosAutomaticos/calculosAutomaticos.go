package calculosAutomaticos

import (
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

var tipoConceptoRemunerativos int = -1
var tipoConceptoNoRemunerativos int = -2
var tipoConceptoDescuentos int = -3

func Hacercalculoautomatico(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {
	var importeCalculado float64
	var importeCalculadoPorPorcentaje float64
	porcentaje := *concepto.Porcentaje / 100
	switch tipocalculo := *concepto.Tipodecalculoid; tipocalculo {
	case -1:
		importeCalculado = calculoRemunerativos(concepto, liquidacion)
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
		return importeCalculadoPorPorcentaje
	case -2:
		importeCalculado = calculoNoRemunerativos(concepto, liquidacion)
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
		return importeCalculadoPorPorcentaje
	case -3:
		importeCalculado = calculoRemunerativosMenosDescuentos(concepto, liquidacion)
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
		return importeCalculadoPorPorcentaje
	case -4:
		importeCalculado = calculoRemunerativosMasNoRemunerativos(concepto, liquidacion)
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
		return importeCalculadoPorPorcentaje
	case -5:
		importeCalculado = calculoRemunerativosMasNoRemunerativosMenosDescuentos(concepto, liquidacion)
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
		return importeCalculadoPorPorcentaje
	default:
		return importeCalculadoPorPorcentaje
	}

}

func calculoRemunerativos(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {
	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoRemunerativos)

	return importeCalculado
}

func calculoNoRemunerativos(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {

	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoNoRemunerativos)
	return importeCalculado
}

func calculoRemunerativosMenosDescuentos(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoRemunerativos)
	importeCalculadoDescuentos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func calculoRemunerativosMasNoRemunerativos(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoNoRemunerativos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos
	return importeCalculado
}

func calculoRemunerativosMasNoRemunerativosMenosDescuentos(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) float64 {

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoNoRemunerativos)
	importeCalculadoDescuentos := calcularImporteSegunTipoConcepto(liquidacion, tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func calcularImporteSegunTipoConcepto(liquidacion *structLiquidacion.Liquidacion, tipoConcepto int) float64 {
	var importeCalculado float64
	var importeNil *float64
	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]

		if *liquidacionitem.Concepto.Tipoconceptoid == tipoConcepto {
			if liquidacionitem.Importeunitario != importeNil {
				importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
			}
		}
	}

	return importeCalculado
}
