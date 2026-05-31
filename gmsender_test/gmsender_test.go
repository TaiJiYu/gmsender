package gmsendertest_test

import (
	"fmt"
	"gmsender/utils"
	"image/color"
	"strings"
	"testing"
)

func TestS(t *testing.T) {
	fmt.Println(utils.DoubleClickTime)

	fmt.Println("[" + strings.Repeat(" ", 3) + "aa" + strings.Repeat(" ", 3) + "]")
}

func TestColor(t *testing.T) {
	a := utils.ColorRGBByOx(0xff0000)
	b := utils.ColorRGBByOx(0x0000ff)

	fmt.Println(a, b)

	var c color.Color
	c = a
	fmt.Println(c)
	c = b
	fmt.Println(c)
	c = a
	fmt.Println(c)
	c = b
	fmt.Println(c)

}
