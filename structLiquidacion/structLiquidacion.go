package structLiquidacion

import (
	"time"

	"github.com/xubiosueldos/conexionBD/structGormModel"
	"github.com/xubiosueldos/legajo/structLegajo"
)

type Liquidacion struct {
	structGormModel.GormModel
	Nombre                               string                  `json:"nombre"`
	Codigo                               string                  `json:"codigo"`
	Descripcion                          string                  `json:"descripcion"`
	Activo                               int                     `json:"activo"`
	Legajo                               *structLegajo.Legajo    `json:"legajo" gorm:"ForeignKey:Legajoid;association_foreignkey:ID;association_autoupdate:false"`
	Legajoid                             *int                    `json:"legajoid" sql:"type:int REFERENCES Legajo(ID)"`
	Tipo                                 *int                    `json:"tipo"`
	Fecha                                time.Time               `json:"fecha"`
	Fechaultimodepositoaportejubilatorio time.Time               `json:"fechaultimodepositoaportejubilatorio"`
	Zonatrabajo                          string                  `json:"zonatrabajo"`
	Condicionpago                        *int                    `json:"condicionpago"`
	Banco                                *int                    `json:"banco"`
	Fechaperiododepositado               time.Time               `json:"fechaperiododepositado"`
	Fechaperiodoliquidacion              time.Time               `json:"fechaperiodoliquidacion"`
	Importesremunerativos                []Importeremunerativo   `json:"importesremunerativos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Importesnoremunerativos              []Importenoremunerativo `json:"importesnoremunerativos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Descuentos                           []Descuento             `json:"descuentos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Retenciones                          []Retencion             `json:"retenciones" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
}
