package Ganancias

type CalculoSACSegundaCuotaExentasNoAlcanzadas struct {
	CalculoGanancias
}

func (cg *CalculoSACSegundaCuotaExentasNoAlcanzadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("SAC_SEGUNDA_CUOTA_EXENTAS_NO_ALCANZADAS")
}

func (cg *CalculoSACSegundaCuotaExentasNoAlcanzadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC_SEGUNDA_CUOTA_EXENTAS_NO_ALCANZADAS", 0, cg)
}

func (cg *CalculoSACSegundaCuotaExentasNoAlcanzadas) getTope() *float64 {
	return nil
}

func (cg *CalculoSACSegundaCuotaExentasNoAlcanzadas) getNombre() string {
	return "SAC Segunda Cuota Exentas/No Alcanzadas"
}

func (cg *CalculoSACSegundaCuotaExentasNoAlcanzadas) getEsMostrable() bool {
	return false
}
