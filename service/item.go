package pezinventory

// Item - object for item data
type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TypeID      string `json:"type-id"`
}

func listItems() (i []Item) {
	i1 := &Item{
		ID:          "1",
		Name:        "Item 1",
		Description: "Item 1 Description",
		TypeID:      "abc",
	}
	i2 := &Item{
		ID:          "2",
		Name:        "Item 2",
		Description: "Item 2 Description",
		TypeID:      "def",
	}
	i = []Item{*i1, *i2}
	return
}

func getItem(id string) (i *Item) {
	i = &Item{
		ID:          id,
		Name:        "Item " + id,
		Description: "Item " + id + " Description",
		TypeID:      "ghi",
	}
	return
}

func listItemsByType(id string) (i []Item) {
	i1 := &Item{
		ID:          "4",
		Name:        "Item 1",
		Description: "Item 1 Description",
		TypeID:      id,
	}
	i2 := &Item{
		ID:          "5",
		Name:        "Item 2",
		Description: "Item 2 Description",
		TypeID:      id,
	}
	i = []Item{*i1, *i2}
	return
}
