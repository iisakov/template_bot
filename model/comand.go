package model

type Comand struct {
	Text string `json:"text"`
}

type ComandsList [][]Comand

type Comands map[int][][]Comand

func (cs Comands) GetComands(stage int) ComandsList {
	return cs[stage]
}

// Comands is Backupable
func (cs Comands) CreateBackup() {}
func (cs Comands) ReadBackup()   {}
