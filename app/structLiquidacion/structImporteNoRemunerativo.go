package structLiquidacion

import "github.com/xubiosueldos/concepto/structConcepto"

type Importenoremunerativo struct {
	Concepto        *structConcepto.Concepto  `json:"concepto" gorm:"ForeignKey:Conceptoid;association_foreignkey:ID;association_autoupdate:false"`
	Conceptoid      *int                      `json:"conceptoid"`
	Cantidad        *int                      `json:"cantidad"`
	Porcentaje      *int                      `json:"porcentaje"`
	Sobreconceptos  []structConcepto.Concepto `json:"sobreconceptos" gorm:"ForeignKey:Sobreconceptoid;association_foreignkey:ID;association_autoupdate:false"`
	Sobreconceptoid *int                      `json:"sobreconceptoid"`
	Importeunitario float32                   `json:"importeunitario" sql:"type:decimal(19,4);"`
	Total           float32                   `json:"total" sql:"type:decimal(19,4);"`
}
