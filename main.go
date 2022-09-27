package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
	}

	values := FindKeyText(string(data))
	fmt.Println(values, len(values))
}

func FindKeyText(data string) []string {
	var v []string

	readValue := func(str string) (string, int) {
		var stack []uint8
		for i := 0; i < len(str); i++ {
			if i > 0 && str[i] == '"' && str[i-1] != '\\' {
				if len(stack) > 0 {
					stack = nil
				} else {
					stack = append(stack, str[i])
				}
			}

			if len(stack) == 0 && (str[i] == ',' || str[i] == '}') {
				return str[:i], i
			}
		}

		return "", 0
	}

	field := "\"text\":"
	fieldLen := len(field)
	for i := 0; i < len(data); {
		if i+fieldLen < len(data) && data[i:i+fieldLen] == field {
			value, n := readValue(data[i+fieldLen:])

			if value != "" {
				value = strings.Trim(value, "\n \t")
				v = append(v, value)
			}
			i = i + fieldLen + n
			continue
		}

		i++
	}
	return v
}
