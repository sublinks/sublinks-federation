package tests

import (
	"testing"
)

/*
This test exists purely to avoid coverage percentage equally NaN
TODO: Delete this file after real tests get added
*/
func Test_true_is_true(t *testing.T) {
	if true != true {
		t.Error("true did not equal true")
	}
}
