package entity

type Tabler interface {
	TableName() string
}

// TableName 会将 User 的表名重写为 `profiles`
func (WaterD) TableName() string {
	return "water_d"
}

//Basedataflow object for REST(CRUD)
type WaterD struct {
	AutoId          uint   `json:"AutoId"`
	MN              uint   `json:"MN"`
	DataTime        string `json:"DataTime"`
	PH              string `json:"PH"`
	WaterLevel      string `json:"WaterLevel"`
	FluorideIon     string `json:"FluorideIon"`
	DissolvedSolids string `json:"DissolvedSolids"`
	Salinity        string `json:"Salinity"`
}
