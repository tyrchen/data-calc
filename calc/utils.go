package calc

import (
	"fmt"
	"math"
)

func Copysign(value, sign Value) Value {
	return Value(math.Copysign(float64(value), float64(sign)))
}

func Abs(v Value) Value {
	return Value(math.Abs(float64(v)))
}

/*
 * increase with 1,3,7 until stop value
 */

func calc137(env *Env) (ret Point) {
	ret.T = false
	last := env.Last
	if last.T {
		ret.V = 1
	} else {
		ret.V = Abs(last.V)*2 + 1

		if ret.V > STOP_VALUE {
			ret.V = 1
		}
	}

	// speical handling in STOP_COL

	if env.Stop && env.CurrentCol >= STOP_COL &&
		(ret.V > STOP_VALUE ||
			last.V == 0 ||
			int(Abs(last.V)) == STOP_VALUE ||
			last.T == true) {
		ret.V = 0
	}
	return
}

func calc137WithSign(env *Env) (ret Point) {
	ret.T = false
	v := Copysign(1, env.Last.V)
	last := env.Last
	if last.T {
		ret.V = v
	} else {
		ret.V = Copysign(Abs(last.V)*2+1, last.V)

		if env.CurrentCol >= STOP_COL &&
			(ret.V > STOP_VALUE ||
				last.V == 0 ||
				last.T == true) {
			ret.V = 0
		}

		if ret.V > STOP_VALUE {
			ret.V = v
		}
	}
	return
}

func calc137forGzf(value Value, upgrade bool) Value {
	if !upgrade {
		return 1
	}
	ret := value*2 + 1
	if ret > 63 {
		return 1
	}
	return ret
}

/*
 * reverse sign if previous.t is true, not reverse otherwise
 */
func calcReverse(env *Env) (ret Point) {
	sign := bsign(env.Bp)
	ret = calc137(env)
	val := Abs(ret.V)
	if env.Last.T {
		ret.V = val * sign * -1
	} else {
		ret.V = val * sign
	}
	return
}

/*
 * follow the sign of up (or reverse)
 */
func calcFollow(env *Env, same Value) (ret Point) {
	sign := bsign(env.Bp)
	ret = calc137(env)
	ret.V *= sign * same
	return ret
}

func signFollow(value Value, bp Bpoint) (ret Value) {
	sign := bsign(bp)
	ret = sign * value
	return ret
}

/*
 * sign of bp: 0: -1, 1: 1
 */
func bsign(bp Bpoint) Value {
	if bp > 0 {
		return 1
	} else {
		return -1
	}
}

func sign(val Point) Value {
	if val.V > 0 {
		return 1
	} else {
		return -1
	}
}

func zsign(val Point) Value {
	if val.T {
		return 1
	} else {
		return -1
	}
}

func pointToBp(p Point) Bpoint {
	if p.V > 0 {
		return 1
	}
	return 0
}

func getNextBp(bp, inst Bpoint) Bpoint {
	return bp ^ ^inst + 2
}

func withZ(p *Point, inst Bpoint) {
	if (p.V > 0 && inst == 1) || (p.V < 0 && inst == 0) {
		p.T = true
	} else {
		p.T = false
	}
}

func calcDelta(data [COLS]Point) (ret Value) {
	var hasz, hasnz Value
	strHasz := "hasz = "
	strHasnz := "hasnz = "
	for i := 0; i < COLS-1; i++ {
		if data[i].T {
			hasz += Abs(data[i].V)
			strHasz += fmt.Sprintf("%d + ", Abs(data[i].V))
		} else {
			hasnz += Abs(data[i].V)
			strHasnz += fmt.Sprintf("%d + ", Abs(data[i].V))
		}
	}
	ret = hasnz - hasz
	fmt.Println(strHasz)
	fmt.Println(strHasnz)
	fmt.Printf("nz: %d, z: %d, ret=%d\n", hasnz, hasz, ret)
	return ret
}

func calcMultiply(data [COLS]Point) (ret bool) {
	if calcDelta(data) > MUL_COND {
		return true
	} else {
		return false
	}
}
