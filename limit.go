package pidentities

import . "github.com/solidifylabs/specops"

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
	const precision = 122
	// First in https://en.wikipedia.org/wiki/List_of_formulae_involving_%CF%80#Iterative_algorithms

	const (
		N = Inverted(DUP1) + iota
		One
		Precision
		BigOne
		A
	)

	const (
		SwapN = Inverted(SWAP1) + iota
		_
		_
		_
		SwapA
	)

	iter := Code{
		Fn(SHR,
			Precision,
			Fn(MUL,
				/* A already on top */
				Fn(ADD,
					BigOne,
					Fn(DIV,
						BigOne,
						Fn(ADD, Fn(SHL, One, N), One),
					),
				),
			),
		),

		Fn(SwapN, Fn(ADD, One, N)), POP,
	}

	return Code{
		PUSH0,           // N
		PC,              // 1
		PUSH(precision), // Precision
		Fn(SHL, Precision, One),
		BigOne, // A

		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,
		iter, iter, iter, iter, iter, iter, iter, iter, iter, iter,

		Fn(SHR, Precision, Fn(MUL, A /* A already on top */)),
		N, SWAP1, DIV,
	}, precision
}
