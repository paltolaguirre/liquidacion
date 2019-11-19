package calculosAutomaticos

import (
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

func Hacercalculoautomatico(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var importeCalculado float64

	switch tipocalculo := *concepto.Tipodecalculoid; tipocalculo {
	case -1:
		importeCalculado = calculoRemunerativos(concepto, liquidacion)
		return importeCalculado
	case -2:
		importeCalculado = calculoNoRemunerativos(concepto, liquidacion)
		return importeCalculado
	case -3:
		importeCalculado = calculoRemunerativosMenosDescuentos(concepto, liquidacion)
		return importeCalculado
	case -4:
		importeCalculado = calculoRemunerativosMasNoRemunerativos(concepto, liquidacion)
		return importeCalculado
	case -5:
		importeCalculado = calculoRemunerativosMasNoRemunerativosMenosDescuentos(concepto, liquidacion)
		return importeCalculado

	default:
		return importeCalculado
	}

}

func calculoRemunerativos(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var tipoConceptoRemunerativos int = -1

	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoRemunerativos)

	return importeCalculado
}

func calculoNoRemunerativos(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var tipoConceptoNoRemunerativos int = -2

	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoNoRemunerativos)
	return importeCalculado
}

func calculoRemunerativosMenosDescuentos(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var tipoConceptoRemunerativos int = -1
	var tipoConceptoDescuentos int = -3

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoRemunerativos)
	importeCalculadoDescuentos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func calculoRemunerativosMasNoRemunerativos(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var tipoConceptoRemunerativos int = -1
	var tipoConceptoNoRemunerativos int = -2

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoNoRemunerativos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos
	return importeCalculado
}

func calculoRemunerativosMasNoRemunerativosMenosDescuentos(concepto structConcepto.Concepto, liquidacion structLiquidacion.Liquidacion) float64 {
	var tipoConceptoRemunerativos int = -1
	var tipoConceptoNoRemunerativos int = -2
	var tipoConceptoDescuentos int = -3

	importeCalculadoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoNoRemunerativos)
	importeCalculadoDescuentos := calcularImporteSegunTipoConcepto(liquidacion, *concepto.Porcentaje, tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func calcularImporteSegunTipoConcepto(liquidacion structLiquidacion.Liquidacion, porcentaje float64, tipoConcepto int) float64 {
	var importeCalculado float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]

		if *liquidacionitem.Concepto.Tipoconceptoid == tipoConcepto {
			importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
		}
	}

	return importeCalculado * porcentaje
}
