package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {

	var In []string
	var Out []int
	var ListOfFiles []int
	maxpoints := 0
	stringsMap := map[string][]int{}

	// Копируем аргументы в отдельный массив
	for i := range os.Args {
		if i == 0 {
			continue
		}
		In = append(In, os.Args[i])
	}

	// Читаем из файла, полученные в прошлой работе данные и записываем из в отображение
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	datastrings := strings.Split(string(data), "\n")
	for i := range datastrings {
		if datastrings[i] == "" {
			break
		}
		keys := strings.Split(datastrings[i], ":")
		values := strings.Split(keys[1], ",")
		for j := range values {
			value, _ := strconv.Atoi(values[j])
			stringsMap[keys[0]] = append(stringsMap[keys[0]], value)
			if ok := contains(ListOfFiles, value); !ok {
				ListOfFiles = append(ListOfFiles, value)
			}
		}
	}
	// Получаем список файлов
	for i := range ListOfFiles {
		count := 0
		for j := range In {
			if ok := contains(stringsMap[In[j]], ListOfFiles[i]); ok {
				count++
			}
		}
		if count > maxpoints {
			maxpoints = count
		}
		Out = append(Out, count)
	}
	// В порядке убывания количества совпадений выводим массив в stdout
	i := maxpoints
	for i != -1 {
		for j := range Out {
			if Out[j] == i {
				fmt.Print(j+1, " file got ", Out[j], " points !\n")
			}
		}
		i--
	}
}

// Имеется ли элемент в списке
func contains(arr []int, element int) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}
