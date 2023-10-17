package cart

import (
	"errors"

	"github.com/ibanezv/go_trafilea_cart/internal/product"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
)

var ErrProductNotFound = errors.New("record not found")
var ErrCartNotFound = errors.New("record not found")

type Carts interface {
	Get(string) (Cart, error)
	Create(string) (Cart, error)
	AddProduct(string, ProductUpdate) (Cart, error)
	ModifyProduct(string, ProductUpdate) (Cart, error)
}

type CartService struct {
	repo repository.Repository
}

func NewCartService(repo repository.Repository) CartService {
	return CartService{repo: repo}
}

func (c *CartService) Get(cartID string) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	return ConvertToCart(cartDb), err
}

func (c *CartService) Create(userID string) (Cart, error) {
	cartDb, err := c.repo.CartCreate(userID)
	return ConvertToCart(cartDb), err
}

func (c *CartService) AddProduct(cartID string, productUpdate ProductUpdate) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	if err != nil {
		return Cart{}, err
	}

	cart := ConvertToCart(cartDb)
	index := findIndexOfProduct(cart.Products, cartID)
	if index >= 0 {
		cart.Products[index].Quantity += productUpdate.Quantity
		return cart, nil
	} else {
		productDB, err := c.repo.ProductGet(productUpdate.ProductID)
		if err != nil {
			return Cart{}, err
		}
		cart.Products = append(cart.Products, product.Product{ProductID: productDB.ProductID, Name: productDB.Name,
			Category: productDB.Category, Quantity: productUpdate.Quantity})
		cartDb, err := c.repo.CartUpdate(convertToDBCart(cart))
		if err != nil {
			return Cart{}, err
		}
		return ConvertToCart(cartDb), nil
	}
}

func (c *CartService) ModifyProduct(cartID string, productUpdate ProductUpdate) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	if err != nil {
		return Cart{}, err
	}

	cart := ConvertToCart(cartDb)
	index := findIndexOfProduct(cart.Products, cartID)
	if index >= 0 {
		cart.Products[index].Quantity += productUpdate.Quantity
		cartDb, err := c.repo.CartUpdate(convertToDBCart(cart))
		if err != nil {
			return Cart{}, err
		}
		return ConvertToCart(cartDb), nil
	} else {
		return Cart{}, ErrCartNotFound
	}
}

func findIndexOfProduct(products []product.Product, productID string) int {
	i := 0
	for i < len(products) {
		if products[i].ProductID == productID {
			return i
		}
		i++
	}
	return -1
}

func ConvertToCart(c repository.CartDB) Cart {
	cart := Cart{CartID: c.CartID, UserID: c.UserID}
	for _, c := range c.Products {
		p := product.Product{ProductID: c.ProductID, Category: c.Category, Name: c.Name, Price: c.Price}
		cart.Products = append(cart.Products, p)
	}
	return cart
}

func convertToDBCart(c Cart) repository.CartDB {
	cartDB := repository.CartDB{CartID: c.CartID, UserID: c.UserID}
	for _, c := range c.Products {
		p := repository.ProductCartDB{ProductDB: repository.ProductDB{ProductID: c.ProductID, Category: c.Category, Name: c.Name, Price: c.Price},
			Quantity: c.Quantity}
		cartDB.Products = append(cartDB.Products, p)
	}
	return cartDB
}