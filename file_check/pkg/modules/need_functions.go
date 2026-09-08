package modules

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func GetHexBytesHash(file multipart.File, w http.ResponseWriter, header *multipart.FileHeader) (fileBytes []byte, fileHash string) {
	fileBytes, err := io.ReadAll(file) // Читаем файл
	if err != nil {
		http.Error(w, "Ошибка чтения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Прочитано %d байт\n", len(fileBytes))

	hashBytes := sha256.Sum256(fileBytes)        // хэш от байт кода
	fileHash = hex.EncodeToString(hashBytes[:]) // переводим в строку и считаем уникальный хекc
	log.Printf("Хэш файла %s: %s\n", header.Filename, fileHash)
	return fileBytes, fileHash
}
