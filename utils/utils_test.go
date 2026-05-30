package utils

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	f := OpenWinChooseFile()
	fmt.Println("内容：", f, filepath.Dir(f))
}
