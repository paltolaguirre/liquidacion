package Ganancias

import (
	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
	"strconv"
)

type CalculoHijos struct {
	CalculoGanancias
}

func (cg *CalculoHijos) getResultInternal() float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig

	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia( "deduccionespersonales", "valorfijomni")
	sql := "SELECT * FROM detallecargofamiliarsiradig WHERE hijoid NOTNULL AND estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64) + "AND detallecargofamiliarsiradig.deleted_at IS NULL"
	cg.Db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {
		porcentaje := detallescargofamiliarsiradig[i].Porcentaje
		importeTotal = importeTotal + cg.getfgDetalleCargoFamiliar("hijoid", "valorfijohijo", *porcentaje)
	}
	return importeTotal
}

func (cg *CalculoHijos) getResult() float64 {
	return cg.getResultOnDemandTemplate("HIJOS", 40, cg)
}

func (cg *CalculoHijos) getTope() *float64 {
	return nil
}

func (cg *CalculoHijos) getNombre() string {
	return "Hijos"
}