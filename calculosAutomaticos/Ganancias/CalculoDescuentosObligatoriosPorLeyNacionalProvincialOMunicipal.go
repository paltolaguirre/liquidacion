package Ganancias

type CalculoDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal struct {
	CalculoGanancias
}

func (cg *CalculoDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("DESCUENTOS_OBLIGATORIOS_POR_LEY_NACIONAL_PROVINCIAL_MUNICIPAL")
}

func (cg *CalculoDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal) getResult() float64 {
	return cg.getResultOnDemandTemplate("Descuentos obligatorios por ley nacional, provincial o municipal (-)", "DESCUENTOS_OBLIGATORIOS_POR_LEY_NACIONAL_PROVINCIAL_MUNICIPAL", 47, cg)
}
