package calculosAutomaticos

import (
	"strconv"
	"time"

	s "strings"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

func getMesesAProrratear(concepto *structConcepto.Concepto, fechaliquidacion *time.Time, db *gorm.DB) int {
	fechadesde := strconv.Itoa(fechaliquidacion.Year()) + "-01-01"
	fechahasta := fechaliquidacion.Format("2006-01-02")
	var fechaliquidacionmasantigua *time.Time
	var sql string

	sql = "SELECT l.fecha FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid WHERE c.id = " + strconv.Itoa(concepto.ID) + " AND l.fecha BETWEEN '" + fechadesde + "' AND '" + fechahasta + "' ORDER BY fecha ASC LIMIT 1"
	fmt.Println("sql: ", sql)
	db.Raw(sql).Row().Scan(&fechaliquidacionmasantigua)

	mesAProrratear := getMes(fechaliquidacionmasantigua)

	return 13 - mesAProrratear

}

func getMes(fecha *time.Time) int {
	mes, _ := strconv.Atoi(s.Split(fecha.String(), "-")[1])
	return mes
}

func getImporteTotalSegunTipoImpuestoGanancias(tipoImpuestoALasGanancias string, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if concepto.Codigo == tipoImpuestoALasGanancias {
			if concepto.Prorrateo == true {
				mes = float64(getMesesAProrratear(concepto, &liquidacion.Fecha, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal
}

func getRemuneracionBruta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", liquidacion, db)
	return importeTotal
}

func getRemuneracionNoHabitual(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES", liquidacion, db)
	return importeTotal
}

/* Consultar lo de los primeros 6 meses*/
func getSacCuotas(liquidacion *structLiquidacion.Liquidacion, correspondeSemestre bool, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if correspondeSemestre && concepto.Basesac == true {
			if concepto.Prorrateo == true {
				mes = float64(getMesesAProrratear(concepto, &liquidacion.Fecha, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal / 12
}

/* Consultar lo de los primeros 6 meses y donde se tendria que tener en cuenta el mes 6*/
func getSacPrimerCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondePrimerSemetre := getMes(&liquidacion.Fecha) < 6
	importeTotal := getSacCuotas(liquidacion, correspondePrimerSemetre, db)

	return importeTotal

}

/* Consultar lo de los segundos 6 meses*/
func getSacSegundaCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondeSegundoSemetre := getMes(&liquidacion.Fecha) < 6
	importeTotal := getSacCuotas(liquidacion, correspondeSegundoSemetre, db)

	return importeTotal
}

func getHorasExtrasGravadas(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getMovilidadYViaticosGravada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getMaterialDidacticoPersonalDocenteRemuneracion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}
