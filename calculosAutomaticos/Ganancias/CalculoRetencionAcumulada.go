package Ganancias

import (
	"strconv"
)

type CalculoRetencionAcumulada struct {
	CalculoGanancias
}

func (cg *CalculoRetencionAcumulada) getResultInternal() float64 {
	anioperiodoliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()
	mesliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
	var totalconceptosimpuestoganancias float64
	var totalconceptosimpuestogananciasdevolucion float64

	sql := "SELECT SUM(li.importeunitario) FROM Liquidacion l INNER JOIN Liquidacionitem li ON l.id = li.liquidacionid INNER JOIN legajo le ON le.id = l.legajoid INNER JOIN concepto c ON c.id = li.conceptoid WHERE to_number(to_char(l.fechaperiodoliquidacion, 'MM'),'99') < " + strconv.Itoa(mesliquidacion) + " AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND le.id = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND c.codigo = 'IMPUESTO_GANANCIAS' AND l.deleted_at IS NULL AND le.deleted_at IS NULL AND li.deleted_at IS NULL AND c.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&totalconceptosimpuestoganancias)

	sql = "SELECT SUM(li.importeunitario) FROM Liquidacion l INNER JOIN Liquidacionitem li ON l.id = li.liquidacionid INNER JOIN legajo le ON le.id = l.legajoid INNER JOIN concepto c ON c.id = li.conceptoid WHERE to_number(to_char(l.fechaperiodoliquidacion, 'MM'),'99') < " + strconv.Itoa(mesliquidacion) + " AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND le.id = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND c.codigo = 'IMPUESTO_GANANCIAS_DEVOLUCION' AND l.deleted_at IS NULL AND le.deleted_at IS NULL AND li.deleted_at IS NULL AND c.deleted_at IS NULL"
	cg.Db.Raw(sql).Row().Scan(&totalconceptosimpuestogananciasdevolucion)

	return totalconceptosimpuestoganancias - totalconceptosimpuestogananciasdevolucion
}

func (cg *CalculoRetencionAcumulada) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETENCION_ACUMULADA", 53, cg)
}

func (cg *CalculoRetencionAcumulada) getTope() *float64 {
	return nil
}

func (cg *CalculoRetencionAcumulada) getNombre() string {
	return "Retención acumulada"
}

func (cg *CalculoRetencionAcumulada) getEsMostrable() bool {
	return true
}
