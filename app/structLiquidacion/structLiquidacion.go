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
	Tipo                                 string                  `json:"tipo"`
	Fecha                                time.Time               `json:"fecha"`
	Fechaultimodepositoaportejubilatorio time.Time               `json:"fechaultimodepositoaportejubilatorio"`
	Zonatrabajo                          string                  `json:"zonatrabajo"`
	Condicionpago                        string                  `json:"condicionpago"`
	Banco                                string                  `json:"banco"`
	Fechainicioperiododepositado         time.Time               `json:"fechainicioperiododepositado"`
	Fechafinperiododepositado            time.Time               `json:"fechafinperiododepositado"`
	Fechainicioperiodoliquidacion        time.Time               `json:"fechainicioperiodoliquidacion"`
	Fechafinperiodoliquidacion           time.Time               `json:"fechafinperiodoliquidacion"`
	Importesremunerativos                []Importeremunerativo   `json:"importesremunerativos" gorm:"ForeignKey:Importeremunerativoid;association_foreignkey:ID;association_autoupdate:false"`
	Importeremunerativoid                *int                    `json:"importeremunerativoid" sql:"type:int REFERENCES Importeremunerativo(ID)"`
	Importesnoremunerativos              []Importenoremunerativo `json:"importesnoremunerativos" gorm:"ForeignKey:Importenoremunerativoid;association_foreignkey:ID;association_autoupdate:false"`
	Importenoremunerativoid              *int                    `json:"importenoremunerativoid" sql:"type:int REFERENCES Importenoremunerativo(ID)"`
	Descuentos                           []Descuento             `json:"descuentos" gorm:"ForeignKey:Descuentoid;association_foreignkey:ID;association_autoupdate:false"`
	Descuentoid                          *int                    `json:"descuentoid" sql:"type:int REFERENCES Descuento(ID)"`
	Retenciones                          []Retencion             `json:"retenciones" gorm:"ForeignKey:Retencionid;association_foreignkey:ID;association_autoupdate:false"`
	Retencionid                          *int                    `json:"retencionid" sql:"type:int REFERENCES Retencion(ID)"`
}
