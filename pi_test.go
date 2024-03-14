package pidentities

import (
	"math"
	"math/big"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/solidifylabs/specops/runopts"
)

func TestApproximations(t *testing.T) {
	tests := []struct {
		name string
		code Implementation
	}{
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
			// Dave
			// It'll annoy him how inefficient the convergence is
			// It also "exhibits unusual behaviour"
			name: "Madhava–Leibniz",
			code: MadhavaLeibniz,
		},
		{
			// Wattsy
			// Takes experimental code to the limit
			name: "Limit",
			code: Limit,
		},
		{
			// Josh
			// Doesn't strike me as the gambling type
			name: "Monte Carlo",
			code: MonteCarlo,
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
	}

	var hashes [][]byte

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := tt.code()

			var (
				startGas uint64
				contract *vm.Contract
			)
			opt := runopts.Func(func(c *runopts.Configuration) error {
				contract = c.Contract
				startGas = contract.Gas
				hashes = append(hashes, crypto.Keccak256(c.Contract.Code))
				return nil
			})

			out, err := code.Run(nil, opt)
			if err != nil {
				t.Fatalf("%T.Run() error %v", code, err)
			}

			gasUsed := int64(startGas - contract.Gas)
			t.Logf("Gas used: %s", humanize.Comma(gasUsed))
			if gasUsed > 25e6 {
				t.Error("Used too much gas; want <25M")
			}
			if want := int64(24_995_000); gasUsed > 40_000 && gasUsed < want {
				t.Errorf("Gas-intensive method not optimised; want at least %s gas used", humanize.Comma(want))
			}

			num, precision := out[:32], out[32:]

			bits := new(big.Int).SetBytes(precision)
			if !bits.IsUint64() {
				t.Fatalf("precision %#x not uint64", precision)
			}

			bigPi := new(big.Rat).SetFrac(
				new(big.Int).SetBytes(num),
				new(big.Int).Lsh(big.NewInt(1), uint(bits.Uint64())),
			)
			t.Log(bigPi.FloatString(76))

			// Drop `exact` flag, not an error
			got, _ := bigPi.Float64()

			if absErr := math.Abs(got - math.Pi); absErr > 0.005 {
				t.Error(got, absErr, math.Log10(absErr))

				// Source: https://gist.github.com/retrohacker/e5fff72b7b75ee058924
				t.Logf("       Pi: 0x3243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c8")
				t.Logf("Numerator: %#x", num)
				t.Logf("     Bits: %d", bits.Uint64())
			}
		})
	}

	t.Logf("Code hashes: %#x", hashes)
}
