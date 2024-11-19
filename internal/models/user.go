package models

type DataReg struct {
	IdUser int32  `json:"id_user" omitempty`
	User   string `json:"user" validate:"required"`
	Mail   string `json:"mail" validate:"required"`
	Tel    string `json:"tel" validate:"required"`
	Pwd    string `json:"pwd" validate:"required"`
}
