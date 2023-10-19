package cart

import (
	"context"
	"errors"

	"github.com/ibanezv/go_trafilea_cart/internal/product"
	"github.com/ibanezv/go_trafilea_cart/internal/repository"
	"github.com/rs/zerolog/log"
)

var ErrProductNotFound = errors.New("product not found")
var ErrCartNotFound = errors.New("cart not found")

type Carts interface {
	Get(context.Context, string) (Cart, error)
	Create(context.Context, string) (Cart, error)
	AddProduct(context.Context, string, ProductUpdate) (Cart, error)
	ModifyProduct(context.Context, string, ProductUpdate) (Cart, error)
}

type CartService struct {
	repo repository.Repository
}

func NewCartService(repo repository.Repository) CartService {
	return CartService{repo: repo}
}

func (c *CartService) Get(ctx context.Context, cartID string) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		log.Printf("cart %s not found:%v", cartID, err)
		return ConvertToCart(cartDb), ErrCartNotFound
	}
	return ConvertToCart(cartDb), err
}

func (c *CartService) Create(ctx context.Context, userID string) (Cart, error) {
	cartDb, err := c.repo.CartCreate(userID)
	return ConvertToCart(cartDb), err
}

func (c *CartService) AddProduct(ctx context.Context, cartID string, productUpdate ProductUpdate) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("cart %s not found:%v", cartID, err)
			return Cart{}, ErrCartNotFound
		}
		return Cart{}, err
	}

	cart := ConvertToCart(cartDb)
	index := findIndexOfProduct(cart.Products, productUpdate.ProductID)
	if index >= 0 {
		cart.Products[index].Quantity += productUpdate.Quantity
	} else {
		productDB, err := c.repo.ProductGet(productUpdate.ProductID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				log.Printf("product %s not found:%v", productUpdate.ProductID, err)
				return Cart{}, ErrProductNotFound
			}
			return Cart{}, err
		}
		cart.Products = append(cart.Products, product.Product{ProductID: productDB.ProductID, Name: productDB.Name,
			Category: productDB.Category, Price: productDB.Price, Quantity: productUpdate.Quantity})
	}

	cartDb, err = c.repo.CartUpdate(convertToDBCart(cart))
	if err != nil {
		return Cart{}, err
	}
	return ConvertToCart(cartDb), nil
}

func (c *CartService) ModifyProduct(ctx context.Context, cartID string, productUpdate ProductUpdate) (Cart, error) {
	cartDb, err := c.repo.CartGet(cartID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("cart %s not found:%v", cartID, err)
			return Cart{}, ErrCartNotFound
		}
		return Cart{}, err
	}

	cart := ConvertToCart(cartDb)
	index := findIndexOfProduct(cart.Products, productUpdate.ProductID)
	if index >= 0 {
		cart.Products[index].Quantity = productUpdate.Quantity
		cartDb, err := c.repo.CartUpdate(convertToDBCart(cart))
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				log.Printf("cart %s not found:%v", cartID, err)
				return Cart{}, ErrCartNotFound
			}
			return Cart{}, err
		}
		return ConvertToCart(cartDb), nil
	} else {
		return Cart{}, ErrProductNotFound
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
		p := product.Product{ProductID: c.ProductID, Category: c.Category, Name: c.Name, Price: c.Price, Quantity: c.Quantity}
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
