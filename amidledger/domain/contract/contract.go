package contract

import (
	"errors"
	"fmt"

	"amidledger/domain/assets"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ShopContract struct {
	contractapi.Contract
}

func (c *ShopContract) Init(ctx contractapi.TransactionContextInterface) error {
	if err := assets.InitShops(ctx); err != nil {
		return err
	}
	if err := assets.InitUsers(ctx); err != nil {
		return err
	}
	return nil
}

func (c *ShopContract) GetUser(ctx contractapi.TransactionContextInterface, login string) (*assets.FabricUser, error) {
	return assets.GetUser(ctx, login)
}

func (c *ShopContract) GetShop(ctx contractapi.TransactionContextInterface, login string) (*assets.ShopUser, error) {
	shop, err := assets.GetShop(ctx, login)
	if err != nil {
		return nil, err
	}
	user, err := assets.GetUser(ctx, login)
	return &assets.ShopUser{User: user, Shop: shop}, err
}

func (c *ShopContract) GetProductList(ctx contractapi.TransactionContextInterface, shopnumber string) ([]*assets.Product, error) {
	shop, err := assets.GetShop(ctx, shopnumber)
	return shop.Products, err
}

func (c *ShopContract) CreateProduct(ctx contractapi.TransactionContextInterface, name string, producer string, date string, safetime int, temp1 int, temp2 int, price float64) error {
	product := assets.Product{Name: name, Producer: producer, Date: date, Safetime: safetime, Temp1: temp1, Temp2: temp2, Amount: 1, Price: price, Shop: "without shop"}
	return assets.SetProduct(ctx, name, &product)
}

func (c *ShopContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*assets.Product, error) {
	return assets.GetAllProducts(ctx)
}

func (c *ShopContract) GetAllShops(ctx contractapi.TransactionContextInterface) ([]*assets.Shop, error) {
	return assets.GetAllShops(ctx)
}

func (c *ShopContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*assets.FabricUser, error) {
	return assets.GetAllUsers(ctx)
}

func (c *ShopContract) Transfer(ctx contractapi.TransactionContextInterface, fromUser string, toUser string, value float64) error {
	userSender, err := assets.GetUser(ctx, fromUser)
	if err != nil {
		return err
	}
	userReceipt, err := assets.GetUser(ctx, toUser)
	if err != nil {
		return err
	}
	if value == 0 {
		return nil
	}
	if userSender.Balance < value {
		return errors.New("balance is lower than value")
	}
	userSender.Balance -= value
	userReceipt.Balance += value
	if err := assets.SetUser(ctx, userSender.Login, userSender); err != nil {
		return err
	}
	if err := assets.SetUser(ctx, userReceipt.Login, userReceipt); err != nil {
		return err
	}
	return err
}

func (c *ShopContract) ShopCalcCost(ctx contractapi.TransactionContextInterface, productName string, amount float64) (float64, error) {
	product, err := assets.GetProduct(ctx, productName)
	if err != nil {
		return 0, nil
	}
	return CalcFinalCost(product.Price, amount), nil
}

func (c *ShopContract) Login(ctx contractapi.TransactionContextInterface, login string, password string) (*assets.FabricUser, error) {
	user, err := assets.GetUser(ctx, login)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("wrong password")
	}
	return user, nil
}

func (c *ShopContract) ShopAcceptCost(ctx contractapi.TransactionContextInterface, shopNumber string, productName string, amount float64) (*assets.ShopTransaction, error) {
	product, err := assets.GetProduct(ctx, productName)
	if err != nil {
		return nil, nil
	}
	shopTxId := fmt.Sprint(assets.RandUint())
	cost, err := c.ShopCalcCost(ctx, productName, amount)
	if err != nil {
		return nil, err
	}
	product.Amount = amount
	product.Shop = shopNumber
	shopTx := &assets.ShopTransaction{Id: shopTxId, ShopNumber: shopNumber, StartCost: cost, CurrentCost: cost, Temp: make([]int, 0), Product: product, TempFailAmount: 0}
	if err := assets.SetShopTx(ctx, shopTx.Id, shopTx); err != nil {
		return nil, err
	}
	c.Transfer(ctx, shopNumber, assets.ProviderLogin, shopTx.StartCost)

	return shopTx, nil
}

func (c *ShopContract) Delivery(ctx contractapi.TransactionContextInterface, shopTxId string) (*assets.ShopTransaction, error) {
	shopTx, err := assets.GetShopTx(ctx, shopTxId)
	if err != nil {
		return nil, err
	}
	temp := assets.RandUintn(100) - 50
	shopTx.Temp = append(shopTx.Temp, temp)
	if temp < shopTx.Product.Temp1 || temp > shopTx.Product.Temp2 {
		shopTx.CurrentCost -= 0.1 * shopTx.StartCost
		shopTx.TempFailAmount += 1
	}
	assets.SetShopTx(ctx, shopTxId, shopTx)
	return shopTx, nil
}

func (c *ShopContract) ShopBuyProduct(ctx contractapi.TransactionContextInterface, shopTxId string, isAccept bool) error {
	shopTx, err := assets.GetShopTx(ctx, shopTxId)
	if err != nil {
		return err
	}
	if !isAccept {
		c.Transfer(ctx, assets.ProviderLogin, shopTx.ShopNumber, shopTx.StartCost)
		return ctx.GetStub().DelState("shoptx@" + shopTxId)
	}
	c.Transfer(ctx, assets.ProviderLogin, shopTx.ShopNumber, shopTx.StartCost-shopTx.CurrentCost)

	shop, err := assets.GetShop(ctx, shopTx.ShopNumber)
	if err != nil {
		return err
	}
	shopTx.Product.Price = shopTx.CurrentCost / shopTx.Product.Amount * 1.5
	shop.Products = append(shop.Products, shopTx.Product)
	if err := ctx.GetStub().DelState("shoptx@" + shopTxId); err != nil {
		return err
	}
	return assets.SetShop(ctx, shop.Number, shop)
}

func (c *ShopContract) GetAllShopTransaction(ctx contractapi.TransactionContextInterface, shopNumber string) ([]*assets.ShopTransaction, error) {
	list, err := assets.GetAllShopTx(ctx)
	txList := make([]*assets.ShopTransaction, 0)
	for _, tx := range list {
		if tx.ShopNumber == shopNumber {
			txList = append(txList, tx)
		}
	}
	return txList, err
}

func (c *ShopContract) BuyProductFromShop(ctx contractapi.TransactionContextInterface, login string, shopNumber string, productName string, amount float64) error {
	user, err := assets.GetUser(ctx, login)
	if err != nil {
		return err
	}
	shop, err := assets.GetShop(ctx, shopNumber)
	if err != nil {
		return err
	}
	shopUser, err := assets.GetUser(ctx, shopNumber)
	if err != nil {
		return err
	}
	for in, pr := range shop.Products {
		if pr.Name != productName {
			continue
		}
		if pr.Amount < amount {
			return errors.New("low amount in shop")
		}
		price := pr.Price * amount
		user.Balance -= price
		shopUser.Balance += price
		product := *pr
		product.Amount = amount
		shop.Products[in].Amount -= amount
		user.Products = append(user.Products, &product)
		if err := assets.SetUser(ctx, user.Login, user); err != nil {
			return err
		}
		if err := assets.SetShop(ctx, shop.Number, shop); err != nil {
			return err
		}
		if err := assets.SetUser(ctx, shopUser.Login, shopUser); err != nil {
			return err
		}
		return nil
	}
	return errors.New("wrong product name")
}

func (c *ShopContract) ReturnProductToShop(ctx contractapi.TransactionContextInterface, login string, shopNumber string, productName string) error {
	productTxId := fmt.Sprint(assets.RandUint())
	productTx := assets.ProductReturnTx{Id: productTxId, ProductName: productName, UserLogin: login, ShopNumber: shopNumber}
	return assets.SetProductTx(ctx, productTxId, &productTx)
}

func (c *ShopContract) GetProductListTxByShop(ctx contractapi.TransactionContextInterface, shopNumber string) ([]*assets.ProductReturnTx, error) {
	prList, err := assets.GetAllProductTx(ctx)
	result := make([]*assets.ProductReturnTx, 0)
	for _, pr := range prList {
		if pr.ShopNumber == shopNumber {
			result = append(result, pr)
		}
	}
	return result, err
}

func (c *ShopContract) AcceptReturnProduct(ctx contractapi.TransactionContextInterface, id string, shopNumber string) error {
	prTx, err := assets.GetProductTx(ctx, id)
	if err != nil {
		return err
	}
	user, err := assets.GetUser(ctx, prTx.UserLogin)
	if err != nil {
		return err
	}
	shop, err := assets.GetShop(ctx, shopNumber)
	if err != nil {
		return err
	}
	shopUser, err := assets.GetUser(ctx, shopNumber)
	if err != nil {
		return err
	}
	for in, pr := range user.Products {
		if pr.Name != prTx.ProductName {
			continue
		}
		price := pr.Amount * pr.Price
		user.Products[in] = user.Products[len(user.Products)-1]
		user.Products = user.Products[:len(user.Products)-1]
		user.Balance += price
		shopUser.Balance -= price
		for index, prod := range shop.Products {
			if prod.Name != pr.Name {
				continue
			}
			shop.Products[index].Amount += pr.Amount
		}
		if err := assets.SetUser(ctx, user.Login, user); err != nil {
			return err
		}
		if err := assets.SetUser(ctx, shopUser.Login, shopUser); err != nil {
			return err
		}
		if err := assets.SetShop(ctx, shopNumber, shop); err != nil {
			return err
		}
		if err := ctx.GetStub().DelState("producttx@" + id); err != nil {
			return err
		}
	}
	return nil
}
