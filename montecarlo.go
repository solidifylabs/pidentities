package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// MonteCarlo estimates pi with entropy sourced from KECCAK256.
func MonteCarlo() Code {
	return convert(monteCarlo)
}

func monteCarlo() (Code, uint8) {
	const bitPrecision = 256 - 32

	// A unit circle inside a 2x2 square covers π/4 of the area. We can
	// (inefficiently) approximate π using sha3 as a source of entropy!
	//
	// Bottom of the stack will always be:
	// - loop total
	// - loops remaining
	// - hit counter (values inside the circle)
	// - constant: 1 (to use DUP instead of PUSH)
	// - constant: 1 << 128 - 1
	// - constant: 1 <<  64 - 1
	// - Entropy (hash)
	//
	// We can therefore use Inverted(DUP/SWAPn) to access them as required,
	// effectively creating variables.
	const (
		Total = Inverted(DUP1) + iota
		Limit
		Hits
		One
		Bits128
		Bits64
		Hash
	)
	const (
		SwapLimit = Limit + 16 + iota
		SwapHits
	)

	return Code{
		PUSH(0x02b000),                         // loop total (~30M gas); kept as the denominator
		DUP1,                                   // loops remaining
		PUSH0,                                  // inside-circle count (numerator)
		PUSH(1),                                // constant-value 1
		Fn(SUB, Fn(SHL, PUSH(0x80), One), One), // 128-bit mask
		Fn(SUB, Fn(SHL, PUSH(0x40), One), One), // 64-bit mask
		stack.ExpectDepth(6),

		JUMPDEST("loop"), stack.SetDepth(6),

		Fn(KECCAK256, PUSH0, PUSH(32)),

		Fn(AND, Bits64, Hash),                    // x = lowest 64 bits
		Fn(AND, Bits64, Fn(SHR, PUSH(64), Hash)), // y = next lowest 64 bits

		Fn(GT,
			Bits128,
			Fn(ADD,
				Fn(MUL, DUP1), // y^2
				SWAP1,         // x^2 <-> y
				Fn(MUL, DUP1), // x^2
			),
		),

		Fn(SwapHits, Fn(ADD, Hits)),

		Fn(JUMPI,
			PUSH("return"),
			Fn(ISZERO, DUP1, Fn(SUB, Limit, One)), // DUP1 uses the top of the stack without consuming it
		),
		stack.ExpectDepth(9),

		SwapLimit, POP, POP,
		Fn(MSTORE, PUSH0),
		Fn(JUMP, PUSH("loop")), stack.ExpectDepth(6),

		JUMPDEST("return"), stack.SetDepth(9),
		// POP, POP,

		Fn(DIV,
			Fn(SHL, PUSH(bitPrecision+2), Hits), // extra 2 to undo π/4
			Total,
		),
	}, bitPrecision
}
