package tests

import (
	"MoTrade/OKXClient"
	"fmt"
	"testing"
)

func TestParamsBuilder(t *testing.T) {
	p := OKXClient.ParamsBuilder().Set("A", "a").Set("D", "d")
	fmt.Println(p)
}
