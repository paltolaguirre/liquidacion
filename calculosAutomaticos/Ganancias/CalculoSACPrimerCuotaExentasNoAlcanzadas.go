package Ganancias

type CalculoSACPrimerCuotaExentasNoAlcanzadas struct {
	CalculoGanancias
}

func (cg *CalculoSACPrimerCuotaExentasNoAlcanzadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("SAC_PRIMER_CUOTA_EXENTAS_NO_ALCANZADAS")
}

func (cg *CalculoSACPrimerCuotaExentasNoAlcanzadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC_PRIMER_CUOTA_EXENTAS_NO_ALCANZADAS", 0, cg)
}

func (cg *CalculoSACPrimerCuotaExentasNoAlcanzadas) getTope() *float64 {
	return nil
}

func (cg *CalculoSACPrimerCuotaExentasNoAlcanzadas) getNombre() string {
	return "SAC Primer Cuota Exentas/No Alcanzadas"
}

func (cg *CalculoSACPrimerCuotaExentasNoAlcanzadas) getEsMostrable() bool {
	return false
}
