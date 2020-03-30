package Ganancias

type CalculoDonacionFisicosNacProvMunArt20 struct {
	CalculoGanancias
}

func (cg *CalculoDonacionFisicosNacProvMunArt20) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "DONACIONES", "deducciondesgravacionsiradig")
	var importeTope float64
	if importeTotal != 0 {
		importeTope = *cg.getTope()
	}
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoDonacionFisicosNacProvMunArt20) getResult() float64 {
	return cg.getResultOnDemandTemplate("DONACION_FISICOS_NAC_PROV_MUN_ART_20", 38, cg)
}

func (cg *CalculoDonacionFisicosNacProvMunArt20) getTope() *float64 {
	importeTope := (&CalculoSubtotal{cg.CalculoGanancias}).getResult() * 0.05 //5% de Subtotal
	return &importeTope
}

func (cg *CalculoDonacionFisicosNacProvMunArt20) getNombre() string {
	return "Donación a los fiscos nac, prov, mun, inst. art. 20 inc. e) y f) LIG (-)"
}

func (cg *CalculoDonacionFisicosNacProvMunArt20) getEsMostrable() bool {
	return true
}
