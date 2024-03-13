package pidentities

import (
	"math"
	"math/big"
	"testing"
)

func TestApproximations(t *testing.T) {
	tests := []struct {
		name string
		code Implementation
	}{
		{
			// Josh
			// Doesn't strike me as the gambling type
			name: "Monte Carlo",
			code: MonteCarlo,
		},
		{
			// Dan
			// It's old-school
			name: "Basel problem",
			code: Basel,
		},
		{
			// Scott
			// Hexapi is leet
			name: "Bailey–Borwein–Plouffe",
			code: BBP,
		},
		{
			// Wattsy
			// Takes experimental code to the limit
			name: "Limit",
			code: Limit,
		},
		{
			// Simon
			// Can't pronounce either of their surnames
			name: "Viète's formula",
			code: Viete,
		},
		{
			// Eugene
			// Wallis was a puzzle solver (https://www.ams.org/journals/bull/1917-24-02/S0002-9904-1917-03015-7/)
			name: "Wallis product",
			code: Wallis,
		},
		{
			// Dave
			// It'll annoy him how inefficient the convergence is
			// It also "exhibits unusual behaviour"
			name: "Madhava–Leibniz",
			code: MadhavaLeibniz,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := tt.code()

			out, err := code.Run(nil)
			if err != nil {
				t.Fatalf("%T.Run() error %v", code, err)
			}

			num, denom := out[:32], out[32:]

			bigPi := new(big.Rat).SetFrac(
				new(big.Int).SetBytes(num),
				new(big.Int).SetBytes(denom),
			)
			t.Log(bigPi.FloatString(76))

			// Drop `exact` flag, not an error
			got, _ := bigPi.Float64()

			if absErr := math.Abs(got - math.Pi); absErr > 0.005 {
				t.Error(got, absErr, math.Log10(absErr))

				// Source: https://gist.github.com/retrohacker/e5fff72b7b75ee058924
				t.Logf("Pi : 0x3243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c8")
				t.Logf("Num: %#x", num)
				t.Logf("Den: %#x", denom)
			}

		})
	}
}
