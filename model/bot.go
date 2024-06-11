package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Backupable interface {
	CreateBackup()
	ReadBackup()
}

func CreateBackup(bv *Backupable, name string) {
	json, err := json.MarshalIndent(bv, "", "	")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile(fmt.Sprintf("backup/backup%s.json", name), json, 0666)
}

func ReadBackup(bv Backupable, name string) {
	if _, err := os.Stat(fmt.Sprintf("backup/backup%s.json", name)); errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	} else {
		f, err := os.ReadFile(fmt.Sprintf("backup/backup%s.json", name))
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(f, bv)
		if err != nil {
			panic(err)
		}
	}
}
