package modules

type Config struct {
	Name string `json:"name"`    // имя модуля
	Enabled bool `json:"enabled"` // включен/выключен
}

type CheckModule struct { // чекаем строку с модулем
	Check []Config `json:"check"`
}
