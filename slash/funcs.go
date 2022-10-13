package slash

func ParseFloat(a any) (float32, bool) {
	if v, k := a.(float32); k {
		return v, true
	}
	if v, k := a.(float64); k {
		return float32(v), true
	}
	return 0, false
}

func ParseString(a any) (string, bool) {
	if v, k := a.(string); k {
		return v, true
	}
	return "", false
}

func ParseInt(a any) (int, bool) {
	f, k := ParseFloat(a)
	if !k {
		return 0, false
	}
	return int(f), true
}
