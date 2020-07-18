package shogi

import "testing"

func TestInv(t *testing.T) {
	tests := []struct {
		input  Color
		expect Color
	}{
		{Black, White},
		{White, Black},
		{NO_COLOR, NO_COLOR},
	}

	for _, test := range tests {
		inv := test.input.Inv()
		if inv != test.expect {
			t.Errorf("%v.Inv(): want %v, got %v\n", test.input, test.expect, inv)
		}
	}
}
