package pidentities

import (
	. "github.com/solidifylabs/specops" //lint:ignore ST1001 SpecOps DSL is designed to be dot-imported
	"github.com/solidifylabs/specops/stack"
)

// MonteCarlo estimates pi with entropy sourced from KECCAK256.
func MonteCarlo() Code {
	return convert(monteCarlo)
}

// monteCarlo is a refactored version of the example from the SpecOps doc.
func monteCarlo() (Code, uint8) {
	const bits = 256 - 32

	// A unit circle inside a 2x2 square covers π/4 of the area. We can
	// (inefficiently) approximate π using sha3 as a source of entropy!
	const (
		total = Inverted(DUP1) + iota
		loopCounter
		hits
		one
		mask128
		mask64
		hash
	)
	const (
		_ = Inverted(SWAP1) + iota
		swapCounter
		swapHits
	)

	return Code{
		PUSH(0x02b000),                         // loop total (~30M gas); kept as the denominator
		total,                                  // loops remaining
		PUSH0,                                  // inside-circle count (numerator)
		PUSH(1),                                // constant-value 1
		Fn(SUB, Fn(SHL, PUSH(0x80), one), one), // 128-bit mask
		Fn(SUB, Fn(SHL, PUSH(0x40), one), one), // 64-bit mask

		stack.ExpectDepth(6),
		JUMPDEST("loop"),
		stack.SetDepth(6),

		Fn(KECCAK256, PUSH0, PUSH(32)),

		Fn(swapHits,
			Fn(ADD,
				hits,
				Fn(GT,
					mask128,
					Fn(ADD,
						Fn(MUL, DUP1, Fn(AND, mask64, hash)),                    // x^2
						Fn(MUL, DUP1, Fn(AND, mask64, Fn(SHR, PUSH(64), hash))), // y^2
					),
				),
			),
		), POP,
		stack.ExpectDepth(7),

		Fn(JUMPI,
			PUSH("return"),
			Fn(ISZERO,
				Fn(swapCounter,
					Fn(SUB, loopCounter, one),
				),
			),
		),
		stack.ExpectDepth(7),

		Fn(MSTORE, PUSH0 /*hash on top*/),
		Fn(JUMP, PUSH("loop")),
		stack.ExpectDepth(6),

		JUMPDEST("return"),
		stack.SetDepth(7),

		Fn(DIV,
			Fn(SHL, PUSH(bits+2), hits), // extra 2 to undo π/4
			total,
		),
	}, bits
}
