package main

import (
	"fmt"
	"strings"
)

func main() {
	original := "/data/texts/{}"
	replacement := "example.txt"

	// 使用strings.Replace替换{}
	result := strings.Replace(original, "{}", replacement, 1)

	fmt.Println(result) // 输出: /data/texts/example.txt
}
