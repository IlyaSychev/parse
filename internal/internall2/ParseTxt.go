package internall2

import (
	"fmt"
	"os"
	"strings"
)

func ParseTxt() {
	// Чтение файла
	data, err := os.ReadFile("NewTXTFile.txt")
	if err != nil {
		fmt.Println("Ошибка чтения txt файла:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var dataArray []string

	for _, line := range lines {
		if strings.Contains(line, "Parameter HGT") {
			dataArray = append(dataArray, line)
		}
	}

	for _, data := range dataArray {
		fmt.Println(data)
	}
}
