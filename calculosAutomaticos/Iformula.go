package calculosAutomaticos

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type iformula interface {
	getResult() float64
}

type RemuneracionBruta struct{
	liquidacion *structLiquidacion.Liquidacion
	db *gorm.DB
}

func (rb *RemuneracionBruta) getResult() float64{
	return getfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", rb.liquidacion, rb.db)
}