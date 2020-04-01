package Ganancias

import "strconv"

type CalculoPagosACuenta struct {
	CalculoGanancias
}

func (cg *CalculoPagosACuenta) getResultInternal() float64 {

	var importeTotal float64
	mesLiquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Format("01")
	anioliquidacion := cg.Liquidacion.Fechaperiodoliquidacion.Year()

	sql := "SELECT SUM(importe) FROM retencionpercepcionsiradig rps INNER JOIN siradigtipogrilla stg ON stg.id = rps.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = rps.siradigid WHERE stg.id IN (-21,-22,-23) AND to_number(to_char(mes, 'MM'),'99') <= " + mesLiquidacion + " AND sdg.legajoid = " + strconv.Itoa(*cg.Liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anioliquidacion) + "' AND rps.deleted_at IS  NULL AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL;"
	cg.Db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func (cg *CalculoPagosACuenta) getResult() float64 {
	return cg.getResultOnDemandTemplate("PAGOS_A_CUENTA", 0, cg)
}

func (cg *CalculoPagosACuenta) getTope() *float64 {
	return nil
}

func (cg *CalculoPagosACuenta) getNombre() string {
	return "Pagos a cuenta"
}

func (cg *CalculoPagosACuenta) getEsMostrable() bool {
	return false
}
