package bank_account_processor

import (
	"testing"
)

func getRequestData() ([][]string, [][]int, [][]int) {
	requestsMap := [][]string{
		{
			"giveMeAHuge 1 3",
		},
		{
			"transfer 1 4 10",
			"deposit 3 10",
			"withdraw 5 15",
		},
		{
			"transfer 1 4 40",
			"deposit 3 10",
			"withdraw 5 65",
		},
		{
			"transfer 2 4 50",
			"transfer 5 8 580",
			"deposit 3 10",
			"deposit 8 10",
			"withdraw 5 65",
			"withdraw 11 65",
			"withdraw 2 10",
		},
		{
			"transfer 1 3 35",
		},
		{
			"transfer 1 3 35",
			"deposit 12 1200",
		},
	}

	balancesMap := [][]int{
		{},
		{20, 30, 10, 90, 60},
		{20, 30, 10, 90, 60},
		{10, 50, 200, 85, 15, 50, 45, 65, 70, 90},
		{45},
		{45, 10, 10},
	}

	exp := [][]int{
		{0}, {10, 30, 20, 100, 45}, {0}, {-1}, {0}, {-1},
	}

	return requestsMap, balancesMap, exp
}

func compareSlices(s1 []int, s2 []int) bool {
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

func TestBankRequests(t *testing.T) {
	requests, balances, exp := getRequestData()
	for i, request := range requests {
		res := bankRequests(request, balances[i])
		equal := compareSlices(exp[i], res)

		if !equal {
			t.Errorf("Failed: Wrong Answer - r: %v i: %v e: %v, a: %v", request, balances[i], exp[i], res)
		}
	}
}
