package model

type Config struct {
	Mysql Mysql `json:"mysql"`
	Redis Redis `json:"redis"`
}

type Mysql struct {
	User string `json:"user"`
	Gorm string `json:"gorm"`
}

type Redis struct {
	Password string `json:"password"`
	Address  string `json:"address"`
}
