package decimal18

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

const PRECISION = 18

var (
	zero = uint256.NewInt(0)

	units = [19]*uint256.Int{
		uint256.NewInt(1e0),
		uint256.NewInt(1e1),
		uint256.NewInt(1e2),
		uint256.NewInt(1e3),
		uint256.NewInt(1e4),
		uint256.NewInt(1e5),
		uint256.NewInt(1e6),
		uint256.NewInt(1e7),
		uint256.NewInt(1e8),
		uint256.NewInt(1e9),
		uint256.NewInt(1e10),
		uint256.NewInt(1e11),
		uint256.NewInt(1e12),
		uint256.NewInt(1e13),
		uint256.NewInt(1e14),
		uint256.NewInt(1e15),
		uint256.NewInt(1e16),
		uint256.NewInt(1e17),
		uint256.NewInt(1e18),
	}
)

func pow10(prec int) *uint256.Int {
	if prec < 0 {
		panic("negative precision")
	}

	if prec <= PRECISION {
		return units[prec]
	}

	return new(uint256.Int).Exp(units[1], uint256.NewInt(uint64(prec)))
}

type Decimal uint256.Int

func NewDecimal(val *uint256.Int, prec int) *Decimal {
	if prec > PRECISION {
		return (*Decimal)(new(uint256.Int).Div(val, pow10(prec-PRECISION)))
	}

	return (*Decimal)(new(uint256.Int).Mul(val, pow10(PRECISION-prec)))
}

func Parse(s string) (*Decimal, error) {
	ss := strings.Split(s, ".")

	switch len(ss) {
	case 1:
		val, err := uint256.FromDecimal(s)
		if err != nil {
			return nil, err
		}
		return NewDecimal(val, 0), nil
	case 2:
		val, err := uint256.FromDecimal(ss[0] + ss[1])
		if err != nil {
			return nil, err
		}
		return NewDecimal(val, len(ss[1])), nil
	default:
		return nil, fmt.Errorf("unsupported fixed point format: %q", s)
	}
}

func (z *Decimal) String() string {
	intPart := new(Decimal).IntPart(z)
	fracPart := new(Decimal).FracPart(z)
	if fracPart.Gt(zero) {
		return intPart.Dec() + "." + fracPart.Dec()
	}

	return intPart.Dec()
}

func (z *Decimal) Significant() *uint256.Int {
	return (*uint256.Int)(z)
}

func (z *Decimal) Shift(x *Decimal, shift int) *Decimal {
	switch {
	case shift == 0:
		return z
	case shift > 0:
		return z.MulInt(x, pow10(shift))
	case shift < 0:
		return z.DivInt(x, pow10(-shift))
	default:
		panic("unreachable")
	}
}

func (z *Decimal) IntPart(x *Decimal) *uint256.Int {
	return z.Significant().Div(x.Significant(), pow10(PRECISION))
}

func (z *Decimal) FracPart(x *Decimal) *uint256.Int {
	return z.Significant().Mod(x.Significant(), pow10(PRECISION))
}

func (z *Decimal) Add(x, y *Decimal) *Decimal {
	return (*Decimal)(z.Significant().Add(x.Significant(), y.Significant()))
}

func (z *Decimal) Sub(x, y *Decimal) *Decimal {
	return (*Decimal)(z.Significant().Sub(x.Significant(), y.Significant()))
}

func (z *Decimal) MulInt(x *Decimal, y *uint256.Int) *Decimal {
	return (*Decimal)(z.Significant().Mul(x.Significant(), y))
}

func (z *Decimal) DivInt(x *Decimal, y *uint256.Int) *Decimal {
	return (*Decimal)(z.Significant().Div(x.Significant(), y))
}

func (z *Decimal) Mul(x, y *Decimal) *Decimal {
	return (*Decimal)(
		z.Significant().Div(
			z.Significant().Mul(x.Significant(), y.Significant()),
			pow10(PRECISION),
		),
	)
}

func (z *Decimal) Div(x, y *Decimal) *Decimal {
	return (*Decimal)(
		z.Significant().Div(
			z.Significant().Mul(x.Significant(), pow10(PRECISION)),
			y.Significant(),
		),
	)
}

func (z *Decimal) Mod(x, y *Decimal) *Decimal {
	return (*Decimal)(z.Significant().Mod(x.Significant(), y.Significant()))
}

func (z *Decimal) Gt(x *Decimal) bool {
	return z.Significant().Gt(x.Significant())
}

func (z *Decimal) Gte(x *Decimal) bool {
	return !z.Significant().Lt(x.Significant())
}

func (z *Decimal) Lt(x *Decimal) bool {
	return z.Significant().Lt(x.Significant())
}

func (z *Decimal) Lte(x *Decimal) bool {
	return !z.Significant().Gt(x.Significant())
}
