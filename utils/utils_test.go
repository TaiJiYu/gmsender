package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	flo := OpenWinFolder()
	fmt.Println("文件夹：", filepath.Join(flo, "hahha.txt"))

	// f := OpenWinSaveFileName()
	// fmt.Println("内容：", f, filepath.Dir(f))
}

func TestRead(t *testing.T) {
	fileS, err := os.Open("test.txt")
	if err != nil {
		// 文件错误
		return
	}
	coon := bytes.NewBuffer(make([]byte, 1024))
	fmt.Println(io.Copy(coon, fileS))
	fmt.Println("done")
}

func TestList(t *testing.T) {
	a := []int{}

	b := []int{1, 2, 3}

	a = make([]int, len(b))
	copy(a, b)
	fmt.Println(a)
	a = a[:0]
	copy(a, b)
	fmt.Println(a)

}

func TestFileName(t *testing.T) {
	s := "C:\\Users\\guangmo\\Documents\\Default.rdp"
	fmt.Println(filepath.Abs(s))
	fmt.Println(filepath.Base(s))
	fmt.Println(filepath.Clean(s))
	fmt.Println(filepath.Dir(s))
	fmt.Println(filepath.EvalSymlinks(s))
	fmt.Println(filepath.Ext(s))
	fmt.Println(filepath.FromSlash(s))
	fmt.Println(filepath.Glob(s))
	fmt.Println(filepath.IsAbs(s))
	fmt.Println(filepath.IsLocal(s))
	fmt.Println(filepath.Localize(s))
	fmt.Println(filepath.Split(s))
	fmt.Println(filepath.ToSlash(s))
	fmt.Println(filepath.VolumeName(s))
}
