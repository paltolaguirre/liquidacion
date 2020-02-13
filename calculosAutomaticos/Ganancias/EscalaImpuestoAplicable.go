package Ganancias

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"strconv"
	s "strings"
)

type strEscalaimpuestoaplicable struct {
	Limiteinferior float64 `json:"limiteinferior"`
	Limitesuperior float64 `json:"limitesuperior"`
	Valorfijo      float64 `json:"valorfijo"`
	Valorvariable  float64 `json:"valorvariable"`
	Mesanio        string  `json:"mesanio"`
}

func getfgEscalaImpuestoAplicable(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) *[]strEscalaimpuestoaplicable {
	var strescalaimpuestoaplicable []strEscalaimpuestoaplicable

	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := s.Split(liquidacion.Fechaperiodoliquidacion.String(), "-")[1]

	mesAnioLiquidacion := mesLiquidacion + "/" + strconv.Itoa(anioLiquidacion)

	sql := "SELECT limiteinferior,limitesuperior,valorfijo,valorvariable,mesanio FROM escalaimpuestoaplicable where mesanio = '" + mesAnioLiquidacion + "' and escalaimpuestoaplicable.deleted_at IS NULL"
	db.Raw(sql).Scan(&strescalaimpuestoaplicable)

	return &strescalaimpuestoaplicable

}
