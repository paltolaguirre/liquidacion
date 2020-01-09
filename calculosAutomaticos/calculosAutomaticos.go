package calculosAutomaticos

import (
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type Calculoautomatico struct {
	Concepto         structConcepto.Concepto       `json:"concepto"`
	Liquidacion      structLiquidacion.Liquidacion `json:"liquidacion"`
	Importecalculado float64                       `json:"importecalculado"`
}

func NewCalculoAutomatico(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion) *Calculoautomatico {
	calculoAutomatico := Calculoautomatico{*concepto, *liquidacion, 0}
	return &calculoAutomatico
}

const tipoConceptoRemunerativos int = -1
const tipoConceptoNoRemunerativos int = -2
const tipoConceptoDescuentos int = -3

func (calculoautomatico *Calculoautomatico) Hacercalculoautomatico() {
	var importeCalculado float64
	var importeCalculadoPorPorcentaje float64
	porcentaje := *calculoautomatico.Concepto.Porcentaje / 100

	switch tipocalculo := *calculoautomatico.Concepto.Tipodecalculoid; tipocalculo {
	case -1:
		importeCalculado = calculoautomatico.calculoRemunerativos()
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje

	case -2:
		importeCalculado = calculoautomatico.calculoNoRemunerativos()
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje

	case -3:
		importeCalculado = calculoautomatico.calculoRemunerativosMenosDescuentos()
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje
	case -4:
		importeCalculado = calculoautomatico.calculoRemunerativosMasNoRemunerativos()
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje

	case -5:
		importeCalculado = calculoautomatico.calculoRemunerativosMasNoRemunerativosMenosDescuentos()
		importeCalculadoPorPorcentaje = importeCalculado * porcentaje

	}
	calculoautomatico.SetImporteCalculado(importeCalculadoPorPorcentaje)
}

func (calculoautomatico *Calculoautomatico) calculoRemunerativos() float64 {
	importeCalculado := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoRemunerativos)

	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) calculoNoRemunerativos() float64 {

	importeCalculado := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoNoRemunerativos)
	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) calculoRemunerativosMenosDescuentos() float64 {

	importeCalculadoRemunerativos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoRemunerativos)
	importeCalculadoDescuentos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) calculoRemunerativosMasNoRemunerativos() float64 {

	importeCalculadoRemunerativos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoNoRemunerativos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos
	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) calculoRemunerativosMasNoRemunerativosMenosDescuentos() float64 {

	importeCalculadoRemunerativos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoRemunerativos)
	importeCalculadoNoRemunerativos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoNoRemunerativos)
	importeCalculadoDescuentos := calculoautomatico.calcularImporteSegunTipoConcepto(tipoConceptoDescuentos)

	importeCalculado := importeCalculadoRemunerativos + importeCalculadoNoRemunerativos - importeCalculadoDescuentos
	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) calcularImporteSegunTipoConcepto(tipoConcepto int) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(calculoautomatico.Liquidacion.Liquidacionitems); i++ {
		liquidacionitem := calculoautomatico.Liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil {
			if *liquidacionitem.Concepto.Tipoconceptoid == tipoConcepto && liquidacionitem.Concepto.ID != calculoautomatico.Concepto.ID {
				if liquidacionitem.Importeunitario != importeNil {
					importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
				}
			}
		}
	}

	return importeCalculado
}

func (calculoautomatico *Calculoautomatico) GetImporteCalculado() float64 {
	return calculoautomatico.Importecalculado
}

func (calculoautomatico *Calculoautomatico) SetImporteCalculado(importe float64) {
	calculoautomatico.Importecalculado = importe
}
