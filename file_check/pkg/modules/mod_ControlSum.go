package modules

import  (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

// передаем байты файла в функцию, где повторно будем расчитывать хэш файла, а потом сравнивать с тем что прислал фронтенд, для проверки, был ли поврежден файл или изменен
func ControlSum_module(fileBytes []byte, hash string) (message string, result bool)  {
	log.Printf("Запуск модуля проверки контрольных сумм.")
	// считаем хэш sha-256 от полученных байт
	hashBytes := sha256.Sum256(fileBytes)
	calcHash := hex.EncodeToString(hashBytes[:])

	log.Printf("Хэш полученного файла: %s\n", calcHash)

	// пока что стоит хэш который расчитывается в Main func, в будущем хэш придет с фронта
	if calcHash != hash { // когда фронт будет отправлять хэш, будет выполняться именно этот блок кода
		return "Контрольная сумма не совпадает! Отклонено", false
	}
	return "Проверка модуля 'Контрольная сумма' завершена. \nУспешно! Целостность файла подтверждена", true
}