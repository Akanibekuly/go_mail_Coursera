package main

import "fmt"

type Payer interface {
	Pay(int) error
}

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("Не достаточно денег в кошельке")
	}
	w.Cash -= amount
	return nil
}

type Card struct {
	Balance    int
	ValidUntil string
	CardHolder string
	CVV        string
	Number     string
}

func (c *Card) Pay(amount int) error {
	if c.Balance < amount {
		return fmt.Errorf("Не хватает денег на карте")
	}
	c.Balance -= amount
	return nil
}

type ApplePay struct {
	Money   int
	AplleId string
}

func (a *ApplePay) Pay(amount int) error {
	if a.Money < amount {
		return fmt.Errorf("Не хватает денег на аккаунте")
	}
	a.Money -= amount
	return nil
}

func Buy(p Payer) {
	switch p.(type) {
	case *Wallet:
		fmt.Println("Оплата наличными?")
	case *Card:
		plasticCard, ok := p.(*Card)
		if !ok {
			fmt.Println("Не удалось проебразовать к типу Card")
		}
		fmt.Println("Вставляйте карту, ", plasticCard.CardHolder)
	default:
		fmt.Println("Что-то новое!")
	}

	err := p.Pay(10)
	if err != nil {
		fmt.Printf("Ошибка при оплате %T: %v\n\n", p, err)
		return
	}
	fmt.Printf("Спасибо за покупку через %T\n\n", p)
}

func main() {
	myWallet := &Wallet{Cash: 100}
	Buy(myWallet)

	var myMoney Payer
	myMoney = &Card{Balance: 100, CardHolder: "Akanibekuly"}
	Buy(myMoney)

	myMoney = &ApplePay{Money: 9}
	Buy(myMoney)
}
