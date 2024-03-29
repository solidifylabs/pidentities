package pidentities

import (
	. "github.com/solidifylabs/specops" //lint:ignore ST1001 SpecOps DSL is designed to be dot-imported
	"github.com/solidifylabs/specops/stack"
)

// Basel implements the Basel problem.
func Basel() Code {
	return convert(basel)
}

func basel() (Code, uint8) {
	const bits = 126

	const (
		one = Inverted(DUP1) + iota
		bigSix
		n
		sum
	)

	const (
		_ = Inverted(SWAP1) + iota
		_
		swapN
	)

	return Code{
		PUSH(1),
		Fn(SHL, PUSH(bits), PUSH(6)),
		one,   // n
		PUSH0, // sum

		JUMPDEST("loop"),
		stack.SetDepth(4),

		Fn(ADD,
			Fn(DIV,
				bigSix,
				Fn(MUL, n, n),
			),
			// sum
		),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(GT,
				PUSH(0x071040), // ~25M gas
				Fn(swapN,
					Fn(ADD, n, one),
				),
			),
		),

		stack.Transform(4)(0),
		Fn(SHL, PUSH(bits)), // sqrt will remove the precision
		sqrt(),
	}, bits
}
