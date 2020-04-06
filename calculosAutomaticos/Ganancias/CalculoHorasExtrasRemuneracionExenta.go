package Ganancias

type CalculoHorasExtrasRemuneracionExenta struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getResultInternal() float64 {
	importeRemExenta := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_EXENTA")
	importeRemGravada := cg.obtenerImporteHorasExtrasCien()
	importeAcumuladorMesAnterior := cg.obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo("HORAS_EXTRAS_REMUNERACION_EXENTA")

	return importeRemExenta + importeRemGravada + importeAcumuladorMesAnterior

}

func (cg *CalculoHorasExtrasRemuneracionExenta) getResult() float64 {
	return cg.getResultOnDemandTemplate("HORAS_EXTRAS_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getTope() *float64 {
	return nil
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getNombre() string {
	return "Horas Extras Remuneración Exenta"
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getEsMostrable() bool {
	return false
}
