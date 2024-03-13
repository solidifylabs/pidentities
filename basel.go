package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// Basel implements the Basel problem.
func Basel() Code {
	return convert(basel)
}

func basel() (Code, uint8) {
	const bitPrecision = 126

	const (
		One = Inverted(DUP1) + iota
		FPSix
		N
		Sum
	)

	const (
		_ = Inverted(SWAP1) + iota
		_
		SwapN
		_
	)

	iter := Code{
		Fn(ADD,
			Fn(DIV,
				FPSix,
				Fn(MUL, N, N),
			),
			// Sum
		),

		Fn(SwapN,
			Fn(ADD, N, One),
		), // Deliberately leaving the old value for a loop counter
	}

	return Code{
		PUSH(1),
		Fn(SHL, PUSH(bitPrecision), PUSH(6)),
		One,   // N
		PUSH0, // Sum

		JUMPDEST("iter"), stack.SetDepth(4),
		iter,
		Fn(JUMPI,
			PUSH("iter"),
			Fn(GT, PUSH(1<<12) /*old N from iter*/),
		),

		stack.Transform(4)(0),
		Fn(SHL, PUSH(bitPrecision)),

		sqrt(),
	}, bitPrecision
}
