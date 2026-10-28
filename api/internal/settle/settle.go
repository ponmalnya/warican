// Package settle は割り勘の精算計算を行う。HTTP や DB には依存しない。
package settle

import (
	"errors"
	"fmt"
)

type Payment struct {
	Payer  string `json:"payer"`
	Amount int    `json:"amount"`
	Memo   string `json:"memo"`
}

// Balance は1人分の過不足。プラスはもらう額、マイナスは払う額。
type Balance struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// Transfer は「From が To に Amount 円渡す」を表す。
type Transfer struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

// Settle はメンバーごとの過不足を members の順で返す。
func Settle(members []string, payments []Payment) ([]Balance, error) {
	if len(members) == 0 {
		return nil, errors.New("メンバーが0人です")
	}

	isMember := map[string]bool{}
	for _, name := range members {
		isMember[name] = true
	}

	paid := map[string]int{}
	total := 0
	for _, p := range payments {
		if !isMember[p.Payer] {
			return nil, fmt.Errorf("%s はメンバーにいません", p.Payer)
		}
		paid[p.Payer] += p.Amount
		total += p.Amount
	}

	// TODO: 割り切れないときの端数をどうするか決める
	perPerson := total / len(members)

	balances := make([]Balance, 0, len(members))
	for _, name := range members {
		balances = append(balances, Balance{Name: name, Amount: paid[name] - perPerson})
	}
	return balances, nil
}

// Transfers は過不足から送金の一覧を作る。払う人と受け取る人を
// balances の順に突き合わせるだけで、送金回数の最小化はしない。
func Transfers(balances []Balance) []Transfer {
	// 残額を書き換えながら進めるので、balances を壊さないようにコピーする。
	rest := make([]Balance, len(balances))
	copy(rest, balances)

	transfers := []Transfer{}
	for i := range rest {
		for j := range rest {
			if rest[i].Amount >= 0 {
				break // i はもう払い終わった
			}
			if rest[j].Amount <= 0 {
				continue // j は受け取る人ではない
			}
			amount := min(-rest[i].Amount, rest[j].Amount)
			transfers = append(transfers, Transfer{From: rest[i].Name, To: rest[j].Name, Amount: amount})
			rest[i].Amount += amount
			rest[j].Amount -= amount
		}
	}
	return transfers
}
