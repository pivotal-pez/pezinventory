package pezinventory

// Type - object for type data
type Type struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"status"`
}

func listTypes() (t []Type) {
	t1 := &Type{
		ID:          "1",
		Name:        "Type 1",
		Description: "Type 1 Description",
		Available:   true,
	}
	t2 := &Type{
		ID:          "2",
		Name:        "Type 2",
		Description: "Type 2 Description",
		Available:   true,
	}
	t = []Type{*t1, *t2}
	return
}

func getType(id string) (t *Type) {
	t = &Type{
		ID:          id,
		Name:        "Type " + id,
		Description: "Type " + id + " Description",
		Available:   true,
	}
	return
}
