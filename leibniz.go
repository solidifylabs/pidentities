package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// MadhavaLeibniz implements the Madhava–Leibniz formulat.
func MadhavaLeibniz() Code {
	return convert(leibniz)
}

func leibniz() (Code, uint8) {
	const precision = 120

	// https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80

	const (
		N = Inverted(DUP1) + iota
		One
		Two
		BigOne
		Result
	)

	const (
		SwapN = Inverted(SWAP1) + iota
	)

	const rounds = 100

	code := Code{
		PUSH(2 * rounds),
		PUSH(1),
		Fn(SHL, One, One), // Two
		Fn(SHL, PUSH(precision), One),
		BigOne, // Result

		stack.ExpectDepth(5),
		JUMPDEST("loop"),
		stack.SetDepth(5),

		Fn(DIV,
			BigOne,
			Fn(ADD, One, Fn(SHL, One, N)),
		),
		ADD,

		Fn(DIV,
			BigOne,
			Fn(ADD, One, Fn(SHL, One, Fn(SUB, N, One))),
		),
		SWAP1, SUB,

		Fn(SwapN,
			Fn(SUB, N, Two),
		),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT, Two),
		),

		Fn(SHL, Two /* top = result */),
	}

	return code, precision
}
