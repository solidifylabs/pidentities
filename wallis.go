package pidentities

import (
	. "github.com/solidifylabs/specops" //lint:ignore ST1001 SpecOps DSL is designed to be dot-imported
	"github.com/solidifylabs/specops/stack"
)

// Wallis implements the Wallis product.
func Wallis() Code {
	return convert(wallis)
}

func wallis() (Code, uint8) {
	const bits = 127

	// https://en.wikipedia.org/wiki/Wallis_product

	const (
		precision = Inverted(DUP1) + iota
		n
		_
		fourNSq
	)

	const (
		_ = Inverted(SWAP1) + iota
		swapN
		swapResult
		swapFourNSq
	)

	code := Code{
		PUSH(bits),
		PUSH(0x049880),              // n (~25M gas)
		Fn(SHL, precision, PUSH(1)), // result

		JUMPDEST("loop"),
		stack.SetDepth(3),

		Fn(SHL,
			PUSH(2),
			Fn(MUL, n, n),
		), // fourNSq

		Fn(SHR,
			precision,
			Fn(MUL,
				Fn(DIV,
					Fn(SHL, precision, swapFourNSq),
					Fn(SUB, fourNSq, PUSH(1)),
				),
				/* top = running product */
			),
		),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT,
				PUSH(1),
				Fn(swapN,
					Fn(SUB, n, PUSH(1)),
				)),
		),

		Fn(SHL, PUSH(1) /* top = result */), // identity is π/2
	}

	return code, bits
}
