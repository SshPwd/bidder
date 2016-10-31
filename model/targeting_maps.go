package model

type Dictionary struct {
	Def  int
	Data map[string]int
}

func (this *Dictionary) Get(key string) int {

	if value, ok := this.Data[key]; ok {
		return value
	}
	return this.Def
}
