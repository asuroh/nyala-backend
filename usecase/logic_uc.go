package usecase

// LogicUC ...
type LogicUC struct {
	*ContractUC
}

// Fibonacci ...
func (uc LogicUC) Fibonacci(N int) (res []float64, err error) {
	n := 1
	for n < N {
		n *= 2
	}

	return res, err
}
