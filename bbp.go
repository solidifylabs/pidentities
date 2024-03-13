package pidentities

import (
	. "github.com/solidifylabs/specops"
)

// BBP implements the Bailey–Borwein–Plouffe formula.
func BBP() Code {
	return convert(bbp)
}

func bbp() (Code, uint8) {
	const precision = 252
	// https://en.wikipedia.org/wiki/Bailey%E2%80%93Borwein%E2%80%93Plouffe_formula

	const (
		SixteenK = Inverted(DUP1) + iota
		EightK
		FourK
		One
		Four
		Five
		Six
		Eight
		BigOne
		BigTwo
		BigFour
		Sum
	)

	const (
		SwapSixteenK = Inverted(SWAP1) + iota
		SwapEightK
		SwapFourK
	)

	fracs := Code{
		Fn(SUB,
			Fn(DIV,
				BigFour,
				Fn(ADD, EightK, One),
			),
			Fn(ADD,
				Fn(DIV,
					BigTwo,
					Fn(ADD, EightK, Four),
				),
				Fn(ADD,
					Fn(DIV,
						BigOne,
						Fn(ADD, EightK, Five),
					),
					Fn(DIV,
						BigOne,
						Fn(ADD, EightK, Six),
					),
				),
			),
		),
	}

	code := Code{
		PUSH0, // 16k
		PUSH0, // 8k
		PUSH0, // 4k
		PUSH(1),
		PUSH(4),
		PUSH(5),
		PUSH(6),
		PUSH(8),
		Fn(SHL, PUSH(precision), One),
		Fn(MUL, BigOne, PUSH(2)),
		Fn(MUL, BigOne, Four),
		PUSH0, // Sum
	}

	body := Code{Fn(SHR, FourK, fracs)}

	for i := 0; i < 32; i++ {
		code = append(
			code,
			Fn(ADD, body /*Sum*/),

			Fn(SwapEightK,
				Fn(ADD, EightK, Eight),
			), POP,
			Fn(SwapSixteenK,
				Fn(SHL, One, EightK),
			), POP,
			Fn(SwapFourK,
				Fn(SHR, One, EightK),
			), POP,
		)
	}

	return code, precision
}
