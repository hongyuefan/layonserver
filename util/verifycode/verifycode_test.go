package verifycode

import (
	"testing"
)

func TestBench(t *testing.T) {

	vid, vul, err := VCodeGenerate(60, 240, 4)

	t.Log(vid, vul, err)

}
