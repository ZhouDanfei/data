package assets

type ShopUser struct {
	User *FabricUser `json:"user"`
	Shop *Shop       `json:"shop"`
}
