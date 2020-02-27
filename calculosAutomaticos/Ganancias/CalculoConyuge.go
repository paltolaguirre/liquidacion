package Ganancias

type CalculoConyuge struct {
	CalculoGanancias
}

func (cg *CalculoConyuge) getResultInternal() float64 {
	return cg.getfgDetalleCargoFamiliar("conyugeid", "valorfijoconyuge", 100)
}

func (cg *CalculoConyuge) getResult() float64 {
	return cg.getResultOnDemandTemplate("CONYUGE", 39, cg)
}

func (cg *CalculoConyuge) getTope() *float64 {
	return nil
}

func (cg *CalculoConyuge) getNombre() string {
	return "Conyuge"
}

func (cg *CalculoConyuge) getEsMostrable() bool {
	return true
}
