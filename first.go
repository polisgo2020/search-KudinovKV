package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	stringsMap := map[string][]int{}
	maxlen := 0

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	j := 0
	for _, f := range files {
		fmt.Println(f.Name())

		data, err := ioutil.ReadFile(os.Args[1] + "//" + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		elems := strings.Split(string(data), " ")
		fmt.Println(elems)
		i := 0
		for {
			if len(elems) == i {
				break
			}
			_, ok := stringsMap[elems[i]]

			if b := contains(stringsMap[elems[i]], j); !ok || !b {
				stringsMap[elems[i]] = append(stringsMap[elems[i]], j)
				if maxlen < len(stringsMap[elems[i]]) {
					maxlen = len(stringsMap[elems[i]])
				}
				fmt.Println("value: ", stringsMap[elems[i]])
			}
			i++
		}
		j++
	}

	outfile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()

	j = maxlen

	for j != 0 {
		for key, value := range stringsMap {
			if len(value) == j {
				var IDs []string

				for _, i := range value {
					IDs = append(IDs, strconv.Itoa(i))
				}

				n, err := outfile.WriteString(key + " : {" + strings.Join(IDs, ", ") + "}\n")

				fmt.Println(n, err)

				if err := outfile.Sync(); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		j--
	}
}

func contains(arr []int, element int) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}
