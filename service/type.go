package pezinventory

type Type struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"status"`
}

func ListTypes() (t []Type) {
	t1 := &Type{
		Id:          "1",
		Name:        "Type 1",
		Description: "Type 1 Description",
		Available:   true,
	}
	t2 := &Type{
		Id:          "2",
		Name:        "Type 2",
		Description: "Type 2 Description",
		Available:   true,
	}
	t = []Type{*t1, *t2}
	return
}

func GetType(id string) (t *Type) {
	t = &Type{
		Id:          id,
		Name:        "Type " + id,
		Description: "Type " + id + " Description",
		Available:   true,
	}
	return
}
