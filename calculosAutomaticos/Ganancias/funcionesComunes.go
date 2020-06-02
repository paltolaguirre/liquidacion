package Ganancias

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"strconv"
	s "strings"
	"time"
)

func getfgMes(fecha *time.Time) int {
	var mes int
	if fecha != nil {
		mes, _ = strconv.Atoi(s.Split(fecha.String(), "-")[1])
	}
	return mes
}

func Sum(s []float64) float64 {
	var sum float64
	for _, val := range s {
		sum += val
	}
	return sum
}

func obtenerTipoImpuesto(concepto *structConcepto.Concepto, db *gorm.DB) string {
	var tipoimpuesto string
	if concepto.Tipoimpuestoganancias != nil {
		tipoimpuesto = concepto.Tipoimpuestoganancias.Codigo
		return tipoimpuesto
	}
	if concepto.Tipoimpuestogananciasid != nil {
		sql := "SELECT codigo FROM tipoimpuestoganancias WHERE id = " + strconv.Itoa(*concepto.Tipoimpuestogananciasid)
		db.Raw(sql).Row().Scan(&tipoimpuesto)
	}

	return tipoimpuesto
}

func getfgImporteTotalTope(importeTotal float64, tope float64) float64 {
	if importeTotal > tope {
		return tope
	} else {
		return importeTotal
	}
}

func obtenerItemGananciaFromLiquidacion(liquidacion *structLiquidacion.Liquidacion) *structLiquidacion.Liquidacionitem {
	var itemGanancia structLiquidacion.Liquidacionitem
	liquidacionItems := liquidacion.Liquidacionitems
	for j := 0; j < len(liquidacionItems); j++ {
		if *liquidacionItems[j].Conceptoid == itemGananciaid || *liquidacionItems[j].Conceptoid == itemGananciaDevolucionid {
			if *liquidacionItems[j].Importeunitario >= 0 {
				itemGanancia = liquidacionItems[j]
				break
			}

		}
	}
	return &itemGanancia
}

const (
	itemGananciaid           = -29
	itemGananciaDevolucionid = -30
)
