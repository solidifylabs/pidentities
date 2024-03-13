package pidentities

import (
	. "github.com/solidifylabs/specops" //lint:ignore ST1001 SpecOps DSL is designed to be dot-imported
	"github.com/solidifylabs/specops/stack"
)

// Limit implements an iterative sequence that, in the limit, approaches a
// function of pi:
//
//	a_0 = 1
//	a_{n+1} = (1 + 1/(2n+1))a_n
//	a^2_n / n -> π
func Limit() Code {
	return convert(limit)
}

func limit() (Code, uint8) {
	const bits = 122

	// First in https://en.wikipedia.org/wiki/List_of_formulae_involving_%CF%80#Iterative_algorithms

	const (
		n = Inverted(DUP1) + iota
		one
		precision
		bigOne
		a
	)
	const (
		SwapN = Inverted(SWAP1) + iota
	)

	return Code{
		PUSH0,                   // n
		PC,                      // 1
		PUSH(bits),              // precision
		Fn(SHL, precision, one), // bigOne
		bigOne,                  // a

		stack.ExpectDepth(5),
		JUMPDEST("loop"),
		stack.SetDepth(5),

		Fn(SHR,
			precision,
			Fn(MUL,
				/* last a already on top */
				Fn(ADD,
					bigOne,
					Fn(DIV,
						bigOne,
						Fn(ADD, Fn(SHL, one, n), one),
					),
				),
			),
		),

		Fn(JUMPI,
			PUSH("loop"),
			Fn(LT,
				Fn(SwapN, Fn(ADD, one, n)),
				PUSH(1000),
			),
		),
		stack.ExpectDepth(5),

		Fn(SHR, precision, Fn(MUL, a /* a already on top */)), // a^2
		Fn(DIV, SWAP1, n),
	}, bits
}
