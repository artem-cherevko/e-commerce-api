package payments

import (
	"e-commerce-api/internal/modules/models"
	"errors"
	"log"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

func getCheckoutSessionLink(userID string, cartID string, cartItems []models.CartItem) (string, error) {
	log.Println("1")
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(cartItems))

	for _, item := range cartItems {
		if item.Quantity <= 0 {
			continue
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String("usd"),
				UnitAmount: stripe.Int64(int64(item.Product.Price * 100)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.Product.Name),
				},
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	if len(lineItems) == 0 {
		return "", errors.New("cart is empty or contains invalid items")
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(stripe.CheckoutSessionModePayment),
		LineItems:  lineItems,
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
	}

	params.AddMetadata("user_id", userID)
	params.AddMetadata("cart_id", cartID)

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
