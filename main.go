package main

import (
	"fmt"
	"party_bot/model"
)

func main() {

	tus := model.Users{}
	tqs := model.Questions{}
	tps := model.Pairs{}

	model.ReadBackup(&tus, "User")
	model.ReadBackup(&tqs, "Question")
	model.ReadBackup(&tps, "Pair")

	tus.AddAnswer("lTest1", "aTest2.2")
	tus.AddAnswer("lTest4", "aTest1.4")

	tus.SetGender("lTest1", 0)
	tus.SetGender("lTest4", 1)

	fmt.Println(tps.FindPair("lTest3"))
	fmt.Println(tps.FindPartner("lTest3"))

}
