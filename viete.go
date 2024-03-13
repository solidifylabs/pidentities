package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// Viete implements Viète's formula.
func Viete() Code {
	return convert(viete)
}

func viete() (Code, uint8) {
	const precision = 127

	// https://en.wikipedia.org/wiki/Vi%C3%A8te%27s_formula

	const (
		Two = Inverted(DUP1) + iota
		Precision
		Result
		BigTwo
		Root
	)

	const (
		_ = Inverted(SWAP1) + iota
		_
		SwapResult
		_
	)

	iter := Code{
		Fn(SHL, Precision,
			Fn(ADD, BigTwo /* top = Root */),
		),

		stack.SetDepth(1),
		sqrtWithCleanup(), // Root
		stack.SetDepth(5),

		Fn(SwapResult,
			Fn(SHR,
				Fn(ADD, Precision, PUSH(1)),
				Fn(MUL,
					Result,
					Root,
				),
			),
		), POP,
	}

	code := Code{
		PUSH(2),
		PUSH(precision),
		Fn(SHL, Precision, PUSH(1)), // Result
		Fn(SHL, Precision, Two),

		PUSH0,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,

		Fn(DIV,
			Fn(SHL, Precision, BigTwo),
			Result,
		),
	}

	return code, precision
}
