package pezinventory

type Item struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TypeId      string `json:"type-id"`
}

func ListItems() (i []Item) {
	i1 := &Item{
		Id:          "1",
		Name:        "Item 1",
		Description: "Item 1 Description",
		TypeId:      "abc",
	}
	i2 := &Item{
		Id:          "2",
		Name:        "Item 2",
		Description: "Item 2 Description",
		TypeId:      "def",
	}
	i = []Item{*i1, *i2}
	return
}

func GetItem(id string) (i *Item) {
	i = &Item{
		Id:          id,
		Name:        "Item " + id,
		Description: "Item " + id + " Description",
		TypeId:      "ghi",
	}
	return
}

func ListItemsByType(id string) (i []Item) {
	i1 := &Item{
		Id:          "4",
		Name:        "Item 1",
		Description: "Item 1 Description",
		TypeId:      id,
	}
	i2 := &Item{
		Id:          "5",
		Name:        "Item 2",
		Description: "Item 2 Description",
		TypeId:      id,
	}
	i = []Item{*i1, *i2}
	return
}
