package modules

import (
	"bytes"
	"log"
)

func FileExtenshion_module(fileBytes []byte) (message string, result bool) {
	log.Println("Запуск модуля проверки расширения файла.")

	fileHeader := fileBytes[:4] // Так как вид файла кроется в первых байтах, берем их

	pdf := []byte{0x25, 0x50, 0x44, 0x46} // %PDF
	exe := []byte{0x4D, 0x5A}             // MZ

	if bytes.Equal(fileHeader[:2], exe) { // в ехе файле важны первые 2 байта
		return "Проверка модуля 'Расширение файла' завершена. \nВаш файл - exe, или возможно какой-то вредоносный скрипт!", false
	}

	if bytes.Equal(fileHeader, pdf) { // в пдф файле важны первые 4 байта
		return "Проверка модуля 'Расширение файла' завершена. \nУспешно! Ваш файл - PDF.", true
	}
	// если файл не пдф и не ехе - то мы выводим его так
	return "Проверка модуля 'Расширение файла' завершена. \nУспешно! Ваш файл не имеет опасного расширения", true
}
