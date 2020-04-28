package csv

import (
	"fmt"
	"testing"
)

func lineHandler(index int, s []string) {

	for _, l := range s {
		fmt.Println(index, ConvertString(l))
	}

}

func TestCsv(t *testing.T) {
	ReadLineCSV("20883020023020220156_20200406_业务明细.csv", lineHandler)
}
