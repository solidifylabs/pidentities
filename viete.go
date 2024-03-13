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
	const bits = 127

	// https://en.wikipedia.org/wiki/Vi%C3%A8te%27s_formula

	const (
		loopCounter = Inverted(DUP1) + iota
		precision
		result
		bigTwo
		root
	)

	const depth = uint(root - loopCounter + 1)

	const (
		swapCounter = Inverted(SWAP1) + iota
		_
		swapResult
	)

	code := Code{
		PUSH(100),                   // loops left
		PUSH(bits),                  // precision
		Fn(SHL, precision, PUSH(1)), // result
		Fn(SHL, precision, PUSH(2)), // constant fixed-precision 2
		PUSH0,                       // root

		stack.ExpectDepth(depth),
		JUMPDEST("loop"),
		stack.SetDepth(depth),

		Fn(SHL, precision, // Shift to prepare for sqrt()
			Fn(ADD, bigTwo /* top = last root */), // The next radicand is always the last root + 2
		),
		stack.ExpectDepth(depth),

		stack.SetDepth(1),     // hack for sqrt() because it expects this; SpecOps needs stack.ExpectDeeperThan().
		sqrtWithCleanup(),     //
		stack.SetDepth(depth), // undo the previous set

		Fn(swapResult,
			Fn(SHR,
				Fn(ADD, precision, PUSH(1)),
				Fn(MUL,
					result,
					root,
				),
			),
		), POP,
		stack.ExpectDepth(depth),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT,
				PUSH0,
				Fn(swapCounter, Fn(SUB, loopCounter, PUSH(1))),
			),
		),
		stack.ExpectDepth(depth),

		Fn(DIV,
			Fn(SHL, precision, bigTwo),
			result,
		),
	}

	return code, bits
}
