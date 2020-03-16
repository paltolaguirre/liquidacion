package Ganancias

import (
	"strconv"

	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
)

type CalculoHijosAnual struct {
	CalculoGanancias
}

func (cg *CalculoHijosAnual) getResultInternal() float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig

	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", "valorfijomni")
	sql := "SELECT * FROM detallecargofamiliarsiradig WHERE hijoid NOTNULL AND estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64) + "AND detallecargofamiliarsiradig.deleted_at IS NULL"
	cg.Db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {
		porcentaje := detallescargofamiliarsiradig[i].Porcentaje
		importeTotal = importeTotal + cg.getfgDetalleCargoFamiliarAnual("hijoid", "valorfijohijo", *porcentaje, valorfijoMNI)
	}
	return importeTotal
}

func (cg *CalculoHijosAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("HIJOS_ANUAL", 0, cg)
}

func (cg *CalculoHijosAnual) getTope() *float64 {
	return nil
}

func (cg *CalculoHijosAnual) getNombre() string {
	return "Hijos Anual"
}

func (cg *CalculoHijosAnual) getEsMostrable() bool {
	return false
}
