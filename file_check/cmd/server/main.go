package main

import (
	"database/sql"
	"encoding/json"
	"file_check/pkg/database"
	"file_check/pkg/modules"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

const port = ":8080" // константа нашего порта локалхоста

func ReadModules() (modules.CheckModule, error) { // смотрим наши модули
	var config modules.CheckModule

	fileBytes, err := os.ReadFile("modules.json") // читаем json
	if err != nil {
		return config, err // возвращаем ошибку если нету файла
	}

	err = json.Unmarshal(fileBytes, &config) // пересобираем наш json в конфиг
	if err != nil {
		return config, err // если ошиька внутри Json'a то возвращаем ошибку
	}

	return config, nil
}

func Upload(w http.ResponseWriter, r *http.Request) { // upload обрабатывает загрузку файла от фронтенда
	if r.Method != http.MethodPost { // Метод запроса только POST
		http.Error(w, "Используйте метод POST.", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // Не дает загрузить больше чем указано байт здесь

	file, header, err := r.FormFile("file") // file — сам файл, header — имя, размер; err — ошибка
	if err != nil {
		http.Error(w, "Не удалось получить файл: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	if header.Size < 5 {
		http.Error(w, "Файл слишком маленький. Загрузите файл боьлше 5 байта", http.StatusBadRequest)
		return
	}

	log.Println("----------------------------------------")
	log.Printf("Получен файл: %s, с размером в %d байт\n", header.Filename, header.Size)

	fileBytes, fileHash := modules.GetHexBytesHash(file, w, header) // я засунул эти преобразования в эту функцию, для более чистого кода, а также возможного изменения

	var status string // переменные для проверки на наличие файла в таблице с таким же хэшем
	var module_info string

	// в слуае если файла нет - будет ошибка, и пропуститься следующий сегмент кода с err == nil, а если файл есть, то ошибки не ббудет
	err = database.DB.QueryRow("SELECT status, module_outputs FROM scan_results WHERE file_hash = ?", fileHash).Scan(&status, &module_info)

	if err == nil { // ошибка в прошлом сегменте есть - файла нет, производятся вычисления, если нету, то берутся данные из таблицы по этому файлу
		w.WriteHeader(http.StatusOK)
		log.Println("Файл был проверен ранее. Беру данные из таблицы без запуска модулей.")
		fmt.Fprintln(w, "----------------------------------------")
		fmt.Fprintf(w, "Файл уже проверен ранее!\nСтатус: %s\nДетали:\n%s\n", status, module_info)
		fmt.Fprint(w, "----------------------------------------")
		return
	}

	if err != sql.ErrNoRows { // ErrNoRows - Та ошибка которую мы ждем когда ищем поле в таблице и НЕ находим. Иная другая - проблема с дб и прерываем работу
		log.Print("Ошибка работы датабазы: ", err)
		http.Error(w, "Ошабка работы датабазы на сервере. Напишите в поддержку,или попробуйте снова", http.StatusInternalServerError)
		return
	}

	cfg, err := ReadModules()
	if err != nil { // проверка данных из прочитанного json файла
		log.Printf("Ошибка запуска модулей")
	}

	var used_modules []string // здесь будет список нашего отчета
	var mod_results []bool   // предсозданные данные для проверок в будущем, установлены обычные значения, на случай если нужные модули отключены
	finalStatus := "SUCCESS. Ваш файл цел и безопасен!"
	message := ""
	mod_result := true

	for _, step := range cfg.Check { // в цикле по нашему конфигу смотрим включены/выключены модули
		if !step.Enabled { // если не включен, то пропускаем
			log.Printf("Модуль '%s' отключен.\n", step.Name)
			continue
		}

		switch step.Name { // делаем поиск по имени файла, в случае неизвестного имени модуля - выводиться предупреждение(ПРОГРАММА НЕ ВЫЛЕТАЕТ)
		case "mod_ControlSum":
			message, mod_result = modules.ControlSum_module(fileBytes, fileHash)
			used_modules = append(used_modules, message)
			mod_results = append(mod_results, mod_result)
		case "mod_FileExtenshion":
			message, mod_result = modules.FileExtenshion_module(fileBytes)
			used_modules = append(used_modules, message)
			mod_results = append(mod_results, mod_result)
		case "mod_Winteon":
			message, mod_result = modules.Winteon_module(header.Filename, fileBytes)
			used_modules = append(used_modules, message)
			mod_results = append(mod_results, mod_result)
		case "mod_AntiVirus":
			result := modules.AV_module()
			used_modules = append(used_modules, result)
		
		default:
			log.Printf("Такого модуля не существует: '%s'", step.Name)
		}
	}
	// и если будет хоть 1 false то будет вывод что файл не прошел проверку
	if slices.Contains(mod_results, false) {
		finalStatus = "FAILED. Ваш файл представляет опасность!"
	}
	modulesReport := strings.Join(used_modules, "\n")

	_, err = database.DB.Exec( // записываем данные в таблицу
		"INSERT INTO scan_results (file_hash, file_name, status, module_outputs) VALUES (?, ?, ?, ?)",
		fileHash, header.Filename, finalStatus, modulesReport,
	)
	if err != nil {
		log.Printf("Результат не сохранен! Ошибка: %v\n", err)

	} else {
		log.Println("Результат сохранен")
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "----------------------------------------") // наш вывод статуса, и информации о каждом модуле который сработал
	fmt.Fprintf(w, "Ваш файл проверен.\nСтатус: %s.\nДетали:\n%s\n", finalStatus, modulesReport)
	fmt.Fprint(w, "----------------------------------------")

	defer file.Close()
}

func main() {

	database.InitDB() // открываем нашу sql таблицу
	
	defer database.DB.Close() // закрываем таблицу в случае падения сервера

	http.HandleFunc("/upload", Upload) // оздаем эндпоинт на отправку

	log.Printf("Сервер запущен на http://localhost" + port) // локал сервер

	http.Handle("/", http.FileServer(http.Dir("./front")))
	err := http.ListenAndServe(port, nil) // проверка запуска
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
