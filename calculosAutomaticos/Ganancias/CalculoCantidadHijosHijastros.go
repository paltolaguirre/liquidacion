package Ganancias

type CalculoCantidadHijosHijastros struct {
	CalculoGanancias
}

func (cg *CalculoCantidadHijosHijastros) getResultInternal() float64 {
	var cantidadHijos float64
	sql := "SELECT count(hijoid) FROM detallecargofamiliarsiradig WHERE hijoid NOTNULL AND estaacargo = true AND detallecargofamiliarsiradig.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&cantidadHijos)
	return cantidadHijos
}

func (cg *CalculoCantidadHijosHijastros) getResult() float64 {
	return cg.getResultOnDemandTemplate("CANTIDAD_HIJOS_HIJASTROS", 0, cg)
}

func (cg *CalculoCantidadHijosHijastros) getTope() *float64 {
	return nil
}

func (cg *CalculoCantidadHijosHijastros) getNombre() string {
	return "Cantidad Hijos/Hijastros"
}

func (cg *CalculoCantidadHijosHijastros) getEsMostrable() bool {
	return false
}
