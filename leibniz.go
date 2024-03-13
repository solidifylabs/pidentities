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
	const bits = 120

	// https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80

	const (
		n = Inverted(DUP1) + iota
		one
		two
		bigOne
		result
	)

	const (
		swapn = Inverted(SWAP1) + iota
	)

	const rounds = 100

	code := Code{
		PUSH(2 * rounds),         // max n (iterates down to zero)
		PUSH(1),                  //
		Fn(SHL, one, one),        // two
		Fn(SHL, PUSH(bits), one), // big one
		bigOne,                   // result

		stack.ExpectDepth(5),
		JUMPDEST("loop"),
		stack.SetDepth(5),

		// n even
		Fn(ADD,
			Fn(DIV,
				bigOne,
				Fn(ADD,
					one,
					Fn(SHL, one, n),
				),
			),
		),

		// n odd
		Fn(SUB,
			SWAP1,
			Fn(DIV,
				bigOne,
				Fn(ADD,
					one,
					Fn(SHL, one, Fn(SUB, n, one)),
				),
			),
		),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT,
				two,
				Fn(swapn,
					Fn(SUB, n, two),
				),
			),
		),

		Fn(SHL, two /* top = result */),
	}

	return code, bits
}
