package contract

func CalcK(amount float64) float64 {
	switch {
	case amount <= 100:
		return 1
	case amount > 100 && amount < 1000:
		return 0.95
	case amount > 1000:
		return 0.9
	default:
		return 1
	}
}

func CalcFinalCost(baseCost float64, amount float64) float64 {
	return amount * baseCost * CalcK(amount)
}
