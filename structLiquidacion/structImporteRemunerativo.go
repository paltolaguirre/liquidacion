package structLiquidacion

import (
	"github.com/xubiosueldos/concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/structGormModel"
)

type Importeremunerativo struct {
	structGormModel.GormModel
	Concepto        *structConcepto.Concepto `json:"concepto" gorm:"ForeignKey:Conceptoid;association_foreignkey:ID;association_autoupdate:false;not null"`
	Conceptoid      *int                     `json:"conceptoid" gorm:"not null"`
	Importeunitario *float32                 `json:"importeunitario" sql:"type:decimal(19,4);" gorm:"not null"`
	Liquidacionid   int                      `json:"liquidacionid"`
	/*
		Cantidad        *int                      `json:"cantidad"`
		Porcentaje      *int                      `json:"porcentaje"`
		Sobreconceptos  []structConcepto.Concepto `json:"sobreconceptos" gorm:"ForeignKey:Sobreconceptoid;association_foreignkey:ID;association_autoupdate:false"`
		Sobreconceptoid *int                      `json:"sobreconceptoid"`
		Total           float32                   `json:"total" sql:"type:decimal(19,4);"`
	*/
}
