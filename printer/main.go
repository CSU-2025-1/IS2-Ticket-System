package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	blacklist = []string{
		"go.mod",
		"go.sum",
		"pkg",
		".env",
		"readme",
		"printer",
		".idea",
		".git",
	}
)

func main() {
	sourceDir := "../../IS2-Ticket-System"
	outputFile := "./out.txt"

	// Открываем итоговый файл для записи
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Рекурсивно обходим директорию
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Открываем файл для чтения
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if slices.Contains(blacklist, file.Name()) {
			return nil
		}

		for i := range blacklist {
			if strings.Contains(path, blacklist[i]) {
				return nil
			}
		}

		// Получаем относительный путь
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Записываем заголовок с путем
		_, err = fmt.Fprintf(outFile, "Файл: %s\n", relPath)
		if err != nil {
			return err
		}

		// Копируем содержимое файла
		_, err = io.Copy(outFile, file)
		if err != nil {
			return err
		}

		// Добавляем пустую строку между файлами
		_, err = fmt.Fprintln(outFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при обработке файлов: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Все файлы успешно объединены в %s\n", outputFile)
}
