package structLiquidacion

import (
	"time"

	"github.com/xubiosueldos/conexionBD/structGormModel"
	"github.com/xubiosueldos/legajo/structLegajo"
)

type Liquidacion struct {
	structGormModel.GormModel
	Nombre                               string                    `json:"nombre"`
	Codigo                               string                    `json:"codigo"`
	Descripcion                          string                    `json:"descripcion"`
	Activo                               int                       `json:"activo"`
	Legajo                               *structLegajo.Legajo      `json:"legajo" gorm:"ForeignKey:Legajoid;association_foreignkey:ID;association_autoupdate:false;auto_preload:true;"`
	Legajoid                             *int                      `json:"legajoid" sql:"type:int REFERENCES Legajo(ID)"`
	Tipo                                 *Liquidaciontipo          `json:"tipo" gorm:"ForeignKey:Tipoid;association_foreignkey:ID;association_autoupdate:false;not null"` /*1-Mensual, 2-Primer Quincena, 3-Segunda Quincena, 4-Vacaciones, 5-SAC, 6-Liquidación Final*/
	Tipoid                               *int                      `json:"tipoid" sql:"type:int REFERENCES Liquidaciontipo(ID)" gorm:"not null"`
	Fecha                                time.Time                 `json:"fecha"`
	Fechaultimodepositoaportejubilatorio time.Time                 `json:"fechaultimodepositoaportejubilatorio"`
	Zonatrabajo                          string                    `json:"zonatrabajo"`
	Condicionpago                        *Liquidacioncondicionpago `json:"condicionpago" gorm:"ForeignKey:Condicionpagoid;association_foreignkey:ID;association_autoupdate:false;not null"` /*1-Cuenta Corriente, 2-Contado*/
	Condicionpagoid                      *int                      `json:"condicionpagoid" sql:"type:int REFERENCES Liquidacioncondicionpago(ID)" gorm:"not null"`
	Cuentabancoid                        *int                      `json:"cuentabancoid" gorm:"not null"`
	//Cuentabanco                          Banco                   `json:"cuentabanco"`
	Fechaperiododepositado     time.Time               `json:"fechaperiododepositado"`
	Fechaperiodoliquidacion    time.Time               `json:"fechaperiodoliquidacion"`
	Importesremunerativos      []Importeremunerativo   `json:"importesremunerativos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Importesnoremunerativos    []Importenoremunerativo `json:"importesnoremunerativos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Descuentos                 []Descuento             `json:"descuentos" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Retenciones                []Retencion             `json:"retenciones" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Aportespatronales          []Aportepatronal        `json:"aportespatronales" gorm:"ForeignKey:Liquidacionid;association_foreignkey:ID"`
	Estacontabilizada          bool                    `json:"estacontabilizada"`
	Asientomanualtransaccionid int                     `json:"asientomanualtransaccionid"`
	Asientomanualnombre        string                  `json:"asientomanualnombre"`
}
