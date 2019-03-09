package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

// GetYangHuiTriangle 获取杨辉三角
func GetYangHuiTriangle(num int) [][]int {
	var result [][]int

	for i := 1; i <= num; i++ {
		var tempArr []int
		for j := 0; j < i; j++ {
			if j == 0 {
				tempArr = append(tempArr, 1)
			} else if i == 1 {
				tempArr = append(tempArr, 1)
			} else if j >= i-1 {
				tempArr = append(tempArr, 1)
			} else {
				tempArr = append(tempArr, result[i-2][j-1]+result[i-2][j])
			}

		}
		result = append(result, tempArr)
	}
	return result
}

// Factorial 求值
func Factorial(n uint64) (result uint64) {
	if n == 0 || n == 1 {
		return 1
	}

	return n * Factorial(n-1)
}

func ReadFile(filepath string) {
	f, error := os.Open(filepath)
	if error != nil {
		fmt.Println("打开文件失败", error)
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.Grow(1024)

	_, err := buf.ReadFrom(f)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println(buf.String())
}

func Processor(seq chan int, wait chan struct{}, layer int) {
	go func() {
		prime, ok := <-seq
		if !ok {
			close(wait)
			return
		}

		fmt.Println("prime ch", layer, prime)
		out := make(chan int)
		Processor(out, wait, layer+1)

		for num := range seq {
			fmt.Println("filter ch", layer, num, prime)
			if num%prime != 0 {
				// fmt.Println("out", num, prime, out)
				out <- num
			}
		}
		close(out)
	}()
}

func braceIndices(s string) ([]int, error) {
	var level, idx int
	var idxs []int

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			level++
			if level == 1 {
				idx = i
			}
		case '}':
			level--
			if level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("unbalanced braces in %q", s)
			}
		}

	}
	if level != 0 {
		return nil, fmt.Errorf("unbalanced braces in %q", s)
	}

	return idxs, nil

}

func parsePathPattern(tpl string) ([]string, string) {
	idxs, err := braceIndices(tpl)
	if err != nil {
		return nil, ""
	}

	keys := make([]string, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')

	for i := 0; i < len(idxs); i++ {
		if i == 0 {
			pattern.WriteString(tpl[0:idxs[i]])
		} else if i%2 != 0 {
			keys[i/2] = tpl[idxs[i-1]:idxs[i]]
			pattern.WriteString(`([a-zA-Z0-9.%]+)`)
			if i+1 < len(idxs) && idxs[i] != idxs[i+1] {
				pattern.WriteString(tpl[idxs[i]:idxs[i+1]])
			} else if i+1 > len(idxs) {
				pattern.WriteString(tpl[idxs[i]:])
			}
		}

	}

	return keys, pattern.String()
}

func main() {
	// 1.  赋值操作
	var user string = "abc"
	fmt.Println("Hello World", user)

	// 2. 循环操作，打印100内的素数
	for i := 1; i < 100; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}

	// 3. 杨辉三角实战
	var nums [][]int = GetYangHuiTriangle(10)
	for i := 0; i < len(nums); i++ {
		fmt.Println(nums[i])
	}

	fmt.Println("Factorial", "15", Factorial(15))

	// 4. 文件IO
	fmt.Println("读取文件")
	ReadFile("/Users/wangys/Downloads/埋点日志.txt")

	// origin, wait := make(chan int), make(chan struct{})
	// Processor(origin, wait, 1)
	// for num := 2; num < 100; num++ {
	// 	fmt.Println("send", num)
	// 	origin <- num
	// }
	// close(origin)
	// <-wait

	// 5. 切片
	var array = [2]string{"hello", "world"}
	var slice1 = array[:]
	slice1[1] = "go"
	slice1 = append(slice1, "test")
	fmt.Printf("%s, %s", array[1], slice1[2])

	// 6. map
	var m = map[string]string{
		"W": "World",
		"X": "XO",
	}
	fmt.Println(m)
	for k, v := range m {
		fmt.Printf("%s=%s\n", k, v)
	}

	// 7. 正则
	var text string
	text = "/food%2d/sdf/test.action"

	reg := regexp.MustCompile(`/([a-zA-Z0-9%]+)/([a-zA-Z0-9%]+)/([a-zA-Z%.]+)`)

	fmt.Printf("%q\n", reg.FindStringSubmatch(text))

	// 8. 匹配
	url := "/getFoodList/{foodCatagory}/{footId}"
	// idxs, _ := braceIndices(url)
	keys, patternStr := parsePathPattern(url)
	fmt.Println("pattern", patternStr)
	for _, v := range keys {
		fmt.Println("key:", v)
	}

	reg = regexp.MustCompile(patternStr)
	fmt.Printf("%q\n", reg.FindStringSubmatch("/getFoodList/1%2D/2sdf.action"))
}
