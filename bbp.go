package pidentities

import (
	. "github.com/solidifylabs/specops" //lint:ignore ST1001 SpecOps DSL is designed to be dot-imported
	"github.com/solidifylabs/specops/stack"
)

// BBP implements the Bailey–Borwein–Plouffe formula.
func BBP() Code {
	return convert(bbp)
}

func bbp() (Code, uint8) {
	const bits = 252

	// https://en.wikipedia.org/wiki/Bailey%E2%80%93Borwein%E2%80%93Plouffe_formula

	const (
		_ = Inverted(DUP1) + iota // 16k
		eightK
		fourK
		one
		four
		five
		six
		eight
		bigone
		bigTwo
		bigFour
		sum
	)

	const (
		swapsixteenK = Inverted(SWAP1) + iota
		swapEightK
		swapFourK
	)

	code := Code{
		PUSH0, // 16k
		PUSH0, // 8k
		PUSH0, // 4k
		PUSH(1),
		PUSH(4),
		PUSH(5),
		PUSH(6),
		PUSH(8),
		Fn(SHL, PUSH(bits), one),
		Fn(MUL, bigone, PUSH(2)),
		Fn(MUL, bigone, four),
		PUSH0, // sum
	}

	fracs := Code{
		Fn(SUB,
			Fn(DIV,
				bigFour,
				Fn(ADD, eightK, one),
			),
			Fn(ADD,
				Fn(DIV,
					bigTwo,
					Fn(ADD, eightK, four),
				),
				Fn(ADD,
					Fn(DIV,
						bigone,
						Fn(ADD, eightK, five),
					),
					Fn(DIV,
						bigone,
						Fn(ADD, eightK, six),
					),
				),
			),
		),
	}

	return append(
		code,

		stack.ExpectDepth(12),
		JUMPDEST("loop"),
		stack.SetDepth(12),

		Fn(ADD,
			Fn(SHR, fourK, fracs),
			/* sum on top*/
		),

		Fn(swapEightK,
			Fn(ADD, eightK, eight),
		), POP,
		Fn(swapsixteenK,
			Fn(SHL, one, eightK),
		), POP,
		Fn(swapFourK,
			Fn(SHR, one, eightK),
		), // Deliberately not popping, to use in loop check

		Fn(JUMPI,
			PUSH("loop"),
			Fn(GT, PUSH(1+bits*4) /* old 4k */),
		),
	), bits
}
