package bank_account_processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var requestScheme = [][]string{
	3: {"action", "from", "amount"},
	4: {"action", "from", "to", "amount"},
}

type BankService interface {
	getActionName() string
	transfer() bool
	withdraw() bool
	deposit() bool
	failedAction()
}

type Bank struct {
	BankService
	action                      string
	requestId, from, to, amount int
	balances                    []int

	requestFailed bool
	failedRequest []int
}

func (bank *Bank) getActionName() string {
	return bank.action
}

func (bank *Bank) transfer() bool {
	invalidAmount := bank.balances[bank.from-1] < bank.amount
	invalidReceiver := len(bank.balances) < bank.to
	if invalidAmount || invalidReceiver {
		return true
	}

	bank.balances[bank.from-1] -= bank.amount
	bank.balances[bank.to-1] += bank.amount
	return false
}

func (bank *Bank) withdraw() bool {
	if bank.balances[bank.from-1] < bank.amount {
		return true
	}

	bank.balances[bank.from-1] -= bank.amount
	return false
}

func (bank *Bank) deposit() bool {
	invalidAccount := len(bank.balances) < bank.from
	if invalidAccount {
		return true
	}

	bank.balances[bank.from-1] += bank.amount
	return false
}

func (bank *Bank) failedAction() {
	failed := []int{-bank.requestId}

	bank.balances = []int{}
	bank.requestFailed = true
	bank.failedRequest = failed
}

func call(bankServ BankService) {
	var requestFailed bool

	action := bankServ.getActionName()
	switch action {
	case "transfer":
		requestFailed = bankServ.transfer()
	case "withdraw":
		requestFailed = bankServ.withdraw()
	case "deposit":
		requestFailed = bankServ.deposit()
	default:
		bankServ.failedAction()
	}

	if requestFailed {
		bankServ.failedAction()
	}
}

func extractRequestParams(request string) map[string]interface{} {
	res := map[string]interface{}{}
	reqSlice := strToSlice(request)
	scheme := requestScheme[len(reqSlice)]

	for i, s := range scheme {
		res[s] = reqSlice[i]
	}

	return res
}

func strToSlice(str string) []interface{} {
	var res []interface{}
	erp := strings.Fields(str)

	for _, v := range erp {
		if isNum, _ := regexp.MatchString("[0-9]+", v); isNum {
			if n, _ := strconv.Atoi(v); n != 0 {
				res = append(res, n)
			}
		} else {
			res = append(res, v)
		}
	}

	return res
}

func bankRequests(requests []string, balances []int) []int {
	var res []int

	for requestId, request := range requests {
		reqParams := extractRequestParams(request)

		bank := &Bank{}
		bank.requestId = requestId
		bank.action = fmt.Sprintf("%s", reqParams["action"])
		bank.amount = reqParams["amount"].(int)
		bank.from = reqParams["from"].(int)

		if _, ok := reqParams["to"]; ok {
			bank.to = reqParams["to"].(int)
		}

		bank.balances = balances
		call(bank)

		if bank.requestFailed {
			return bank.failedRequest
		}
		res = bank.balances
	}

	return res
}
