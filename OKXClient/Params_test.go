package OKXClient

import (
	"fmt"
	"testing"
)

func TestParamsBuilder(t *testing.T) {
	p := ParamsBuilder().Set("A", "a").Set("D", "d")
	fmt.Println(p)
}
