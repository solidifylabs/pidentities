// Package pidentities implements ridiculous ways to approximate π through
// hand-written EVM bytecode.
package pidentities

import (
	. "github.com/solidifylabs/specops"
	"github.com/solidifylabs/specops/stack"
)

// An Implementation implements a pidentity.
type Implementation func() Code

// An implementation is a raw implementation that leaves pi on the top of the
// stack.
type implementation func() (_ Code, bitsOfPrecision uint8)

// convert returns an exportable Implementation from a raw implementation. It
// uses the precision to calculate a denominator, treating the value on the top
// of the stack as a numerator.
func convert(fn implementation) Code {
	code, bits := fn()
	return Code{
		code,

		Fn(MSTORE, PUSH0 /*pi on top of the stack*/),
		Fn(MSTORE,
			PUSH(32),
			Fn(SHL, PUSH(bits), PUSH(1)),
		),
		Fn(RETURN, PUSH0, PUSH(64)),
	}
}

// sqrtWithCleanup is the same as sqrt() except it removes its scratchpad from
// the stack.
func sqrtWithCleanup() Code {
	return Code{
		sqrt(),
		SWAP4,
		POP, POP, POP, POP, SWAP1, POP,
		stack.ExpectDepth(1),
	}
}

// sqrt is the square-root function from the SpecOps examples. It consumes the
// value at the top of the stack, leaving its square root instead. It also adds
// other values used in the calculation (use sqrtWithCleanup() to remove these).
func sqrt() Code {
	const (
		Input = Inverted(DUP1) + iota
		One
		ThresholdBits
		Threshold
		xAux
		Result
		Branch
	)
	const (
		SwapInput = Input + 16 + iota
		_         // SetOne
		SetThresholdBits
		SetThreshold
		SetXAux
		SetResult
		SetBranch
	)

	approx := Code{
		stack.ExpectDepth(6),
		Fn(GT, xAux, Threshold), // Branch

		Fn(SetXAux,
			Fn(SHR,
				Fn(MUL, ThresholdBits, Branch),
				xAux,
			),
		), POP, // old value

		Fn(SetThresholdBits,
			Fn(SHR, One, ThresholdBits),
		), POP,

		Fn(SetThreshold,
			Fn(SUB, Fn(SHL, ThresholdBits, One), One),
		), POP,

		Fn(SetResult,
			Fn(SHL,
				Fn(MUL, ThresholdBits, Branch),
				Result,
			),
		), POP,

		POP, // Branch
		stack.ExpectDepth(6),
	}

	// Single round of Newton–Raphson
	newton := Code{
		stack.ExpectDepth(6),
		Fn(SetResult,
			Fn(SHR,
				One,
				Fn(ADD,
					Result,
					Fn(DIV, Input, Result),
				),
			),
		), POP,
		stack.ExpectDepth(6),
	}

	return Code{
		stack.ExpectDepth(1), // Input
		PUSH(1),              // One
		PUSH(128),            // ThresholdBits
		Fn(SUB, Fn(SHL, ThresholdBits, One), One), // Threshold
		Input, // xAux := Input
		One,   // Result
		stack.ExpectDepth(6),

		approx, approx, approx, approx, approx, approx, approx,
		stack.ExpectDepth(6),
		newton, newton, newton, newton, newton, newton, newton,
	}
}
