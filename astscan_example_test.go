package astscan

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"testing"
)

var (
	src = `// main 测试包

package main

import "fmt"

var (
	hello = "你好，世界"  // hello 的注释
)

func main() {
	fmt.Println(hello)
	fmt.Println("测试 2333")
}`
)

func _print(pkg, fileline, value string, n ast.Node) {
	fmt.Printf("pkg:%s, fileline:%s, value:%s\n", pkg, fileline, value)
}

func TestFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(src)); err != nil {
		t.Fatal(err)
	}

	if err := File(tmpFile.Name(), ContainChinese, _print); err != nil {
		t.Fatal(err)
	}
}

func TestDir(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(src)); err != nil {
		t.Fatal(err)
	}

	if err := Dir(os.TempDir(), ContainChinese, _print); err != nil {
		t.Fatal(err)
	}
}
