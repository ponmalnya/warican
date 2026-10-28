package main

import (
	"errors"
	"fmt"
	"os"
)

type Payment struct {
	Payer  string
	Amount int
	Memo   string
}

// Settle はメンバーごとの過不足を返す。プラスはもらう額、マイナスは払う額。
// 計算だけを担当し、表示はしない。
func Settle(members []string, payments []Payment) (map[string]int, error) {
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

	diffs := map[string]int{}
	for _, name := range members {
		diffs[name] = paid[name] - perPerson
	}
	return diffs, nil
}

// Transfer は「From が To に Amount 円渡す」を表す。
type Transfer struct {
	From   string
	To     string
	Amount int
}

// Transfers は過不足から送金の一覧を作る。払う人と受け取る人を
// members の順に突き合わせるだけで、送金回数の最小化はしない。
func Transfers(members []string, diffs map[string]int) []Transfer {
	// 残額を書き換えながら進めるので、diffs を壊さないようにコピーする。
	rest := map[string]int{}
	for name, d := range diffs {
		rest[name] = d
	}

	var transfers []Transfer
	for _, from := range members {
		for _, to := range members {
			if rest[from] >= 0 {
				break // from はもう払い終わった
			}
			if rest[to] <= 0 {
				continue // to は受け取る人ではない
			}
			amount := min(-rest[from], rest[to])
			transfers = append(transfers, Transfer{From: from, To: to, Amount: amount})
			rest[from] += amount
			rest[to] -= amount
		}
	}
	return transfers
}

func main() {
	members := []string{"三郎", "はなこ", "じろう"}

	payments := []Payment{
		{Payer: "三郎", Amount: 0, Memo: "宿代"},
		{Payer: "はなこ", Amount: 0, Memo: "昼ごはん"},
		{Payer: "じろう", Amount: 1500, Memo: "電車"},
	}

	diffs, err := Settle(members, payments)
	if err != nil {
		fmt.Println("エラー:", err)
		os.Exit(1)
	}

	for _, name := range members {
		diff := diffs[name]
		if diff > 0 {
			fmt.Printf("%s: %d円 もらう\n", name, diff)
		} else if diff < 0 {
			fmt.Printf("%s: %d円 払う\n", name, -diff)
		} else {
			fmt.Printf("%s: 精算なし\n", name)
		}
	}

	fmt.Println("--- 送金 ---")
	for _, t := range Transfers(members, diffs) {
		fmt.Printf("%s → %s: %d円\n", t.From, t.To, t.Amount)
	}
}
