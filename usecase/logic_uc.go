package usecase

import (
	"nyala-backend/pkg/str"
	"nyala-backend/usecase/viewmodel"
)

// LogicUC ...
type LogicUC struct {
	*ContractUC
}

// Fibonacci ...
func (uc LogicUC) Fibonacci(N int) (res []int, err error) {
	a := 0
	b := 1
	for a <= N {
		res = append(res, a)
		temp := a
		a = b
		b = temp + a
	}
	return res, err
}

// Prima ..
func (uc LogicUC) Prima(N int) (res []int, err error) {
	for i := 1; i <= N; i++ {
		bil := 0
		for j := 1; j <= i; j++ {
			if i%j == 0 {
				bil = bil + 1
			}
		}
		if bil == 2 {
			res = append(res, i)
		}
	}
	return res, err
}

// CheckPalindrome ...
func (uc LogicUC) CheckPalindrome(kata string) (res viewmodel.LogicPalindromeVM, err error) {
	res.KataPalindrome = str.Reverse(kata)
	res.KataAsli = kata
	if res.KataAsli == res.KataPalindrome {
		res.IsPalindrome = true
	} else {
		res.IsPalindrome = false
	}

	return res, err
}
