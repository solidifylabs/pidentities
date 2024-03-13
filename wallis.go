package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// Wallis implements the Wallis product.
func Wallis() Code {
	return convert(wallis)
}

func wallis() (Code, uint8) {
	const precision = 127
	// https://en.wikipedia.org/wiki/Wallis_product

	const (
		Precision = Inverted(DUP1) + iota
		N
		Result
		FourNSq
	)

	const (
		_ = Inverted(SWAP1) + iota
		SwapN
		SwapResult
		SwapFourNSq
	)

	code := Code{
		PUSH(precision),
		PUSH(200),                   // N,
		Fn(SHL, Precision, PUSH(1)), // Result

		JUMPDEST("loop"), stack.SetDepth(3),

		Fn(SHL,
			PUSH(2),
			Fn(MUL, N, N),
		), // FourNSq

		Fn(DIV,
			Fn(SHL, Precision, SwapFourNSq),
			Fn(SUB, FourNSq, PUSH(1)),
		),

		Fn(SHR,
			Precision,
			Fn(MUL /* top 2 are Result and fraction */),
		),

		Fn(SwapN,
			Fn(SUB, N, PUSH(1)),
		),
		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT, PUSH(1) /* top = last loop counter*/),
		),

		Fn(SHL, PUSH(1) /* top = result */),
	}

	return code, precision
}
