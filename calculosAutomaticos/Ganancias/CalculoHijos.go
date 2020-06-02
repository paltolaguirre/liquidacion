package Ganancias

import (
	"strconv"

	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
)

type CalculoHijos struct {
	CalculoGanancias
}

func (cg *CalculoHijos) getResultInternal() float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig
	var porcentaje float64 = 0
	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", "valorfijomni")
	sql := "SELECT * FROM detallecargofamiliarsiradig inner join hijo on hijo.id = detallecargofamiliarsiradig.hijoid WHERE estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64) + "AND detallecargofamiliarsiradig.deleted_at IS NULL AND hijo.legajoid =" + strconv.Itoa(*cg.Liquidacion.Legajoid)
	cg.Db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {

		if detallescargofamiliarsiradig[i].Porcentaje != nil {
			porcentaje = *detallescargofamiliarsiradig[i].Porcentaje
		}
		importeTotal = importeTotal + cg.getfgDetalleCargoFamiliar("hijoid", "valorfijohijo", porcentaje)
	}
	return importeTotal
}

func (cg *CalculoHijos) getResult() float64 {
	return cg.getResultOnDemandTemplate("HIJOS", 41, cg)
}

func (cg *CalculoHijos) getTope() *float64 {
	return nil
}

func (cg *CalculoHijos) getNombre() string {
	return "Hijos"
}

func (cg *CalculoHijos) getEsMostrable() bool {
	return true
}
