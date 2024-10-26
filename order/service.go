package order

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("request canceled or timed out")
	default:
	}

	if err := validateProducts(products); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, errors.New("no products ordered")
	}

	totalPrice := calculateTotalPrice(products)
	o := &Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		Products:   products,
		TotalPrice: totalPrice,
	}

	err := s.repository.PutOrder(ctx, *o)
	if err != nil {
		log.Printf("failed to post order for account %s: %v", accountID, err)
		return nil, errors.Wrap(err, "failed to store order in repository")
	}

	return o, nil
}

func calculateTotalPrice(products []OrderedProduct) float64 {
	totalPrice := decimal.NewFromFloat(0.0)
	for _, p := range products {
		productPrice := decimal.NewFromFloat(p.Price).Mul(decimal.NewFromFloat(float64(p.Quantity)))
		totalPrice = totalPrice.Add(productPrice)
	}

	return totalPrice.InexactFloat64()
}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}

func validateProducts(products []OrderedProduct) error {
	for _, p := range products {
		if p.ID == "" && p.Name == "" {
			return errors.New("product ID and Name cannot be empty")
		}
		if p.Price < 0 {
			return errors.New("product price cannot be negative")
		}
		if p.Quantity <= 0 {
			return errors.New("product quantity must be greater than zero")
		}
	}

	return nil
}
