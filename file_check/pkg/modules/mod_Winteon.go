package modules

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
)

type OurRequest struct { // формат в котором отправим файл
	SendTime  string `json:"sent_at"`
	FileName  string `json:"file_name"`
	FileBytes string `json:"file_bytes"`
}
type GetTheyResult struct { // то что мы примем от Winteon
	Status        string `json:"status"`
	Message string `json:"message"`
}

func WaitingTimer(timeout, ticker <-chan time.Time) (JsonVariant string) { // функция для нашего таймера
	for {
		select {
		case <-timeout:
			return "Время вышло, ответа нет"
		case <-ticker:
			return "json получен"
		}
	}
}

func Winteon_module(fileName string, fileBytes []byte) (message string, result bool) { // принимаем название файла и его данные
	log.Println("Создаем Json Для отправки файла на Winteon")

	base64File := base64.StdEncoding.EncodeToString(fileBytes) // кодируем данные в формат базы64 для отправки json'ом

	payload := OurRequest{ // наш json
		SendTime: time.Now().Format(time.RFC3339), // задаем время создания файла как текущее(на момент создания)
		FileName:  fileName,
		FileBytes:  base64File,
	}

	_, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "Json не создан: " + err.Error(), false
	}
	
	log.Printf("Json создан и готов к отправке!")
	// здесь должен быть апи с отправкой файла в программе сверху (для этого у сегмента кода выше, вместо _ должна быть переменная, и она должна отправится)
	// предположим сейчас снизу будет небольшая часть кода которая просто создаст json ответа на наш файл и выведет нам ответ

	timeout := time.After(15 * time.Second)
	ticker := time.Tick(3 * time.Second)
	JsonVariant := WaitingTimer(timeout, ticker) // по сути тут мы должны получать файл, и он будет идти ниже

	// здесь я перезаписываю переменную, потому что неоткуда пока что получать нужный json
	JsonVariant = `{ 
		"status": "success",
		"message": "Файл чист и без аномалий!"
	}`

	var apiResult GetTheyResult // создаем перменную с нашей структурой
	err = json.Unmarshal([]byte(JsonVariant), &apiResult)
	if err != nil {
		return "Ответ не пришел: " + err.Error(), false
	} // если не поулчилось - ошибка

	if apiResult.Status == "success" { // проверка статуса, если все хорошо - вывод ответа, если нет, пишем что статус неудовлетворительный
		return "Проверка модуля 'Winteon' завершена. \nУспешно! Результат: " + apiResult.Message, true
	}
	return "Проверка модуля 'Winteon' завершена с статусом: " + apiResult.Status, false
}
