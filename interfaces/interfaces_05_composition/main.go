package main

import "fmt"

type Phone struct {
	Money   int
	AppleId string
}

func (p *Phone) Pay(amount int) error {
	if p.Money < amount {
		return fmt.Errorf("Not enough money on account")
	}
	p.Money -= amount
	return nil
}

func (p *Phone) Ring(number string) error {
	if number == "" {
		return fmt.Errorf("Please, enter phone")
	}
	return nil
}

type Payer interface {
	Pay(int) error
}

type Ringer interface {
	Ring(string) error
}

type NFCPhone interface {
	Payer
	Ringer
}

func PayForMetwiWithMoney(phone NFCPhone) {
	err := phone.Pay(1)
	if err != nil {
		fmt.Printf("Ошибка при оплате %v\n\n", err)
		return
	}
	fmt.Printf("Турникет открыт через %T\n\n", phone)
}

func main() {
	myPhone := &Phone{Money: 9}
	PayForMetwiWithMoney(myPhone)
}
