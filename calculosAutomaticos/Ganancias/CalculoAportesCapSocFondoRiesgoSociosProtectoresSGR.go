package Ganancias

type CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR struct {
	CalculoGanancias
}

func (cg *CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("Montoreintegrar + Montoreintegrar3", "REINTEGRO_DE_APORTES_DE_SOCIOS_PROTECTORES_A_SOCIEDADES_DE_GARANTIA_RECIPROCA", "ajustesiradig")
}

func (cg *CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR) getResult() float64 {
	return cg.getResultOnDemandTemplate("Aportes cap. Soc. / Fondo de riesgo de socios protectores de SGR (-)", "APORTES_CAP_SOC_FONDO_RIESGO_SOCIOS_PROTECTORES_SGR", 30, cg)
}

func (cg *CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR) getTope() *float64 {
	return nil
}
