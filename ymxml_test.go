package YandexMarketXml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestYandexMarketYML_Parse(t *testing.T) {
	v, err := YandexMarketXML("testdata/vendorModel.xml").Parse()
	if err != nil {
		fmt.Printf("%v", err)
	}
	assert.Equal(t, v.Shop.Name, "BestSeller")
	assert.Equal(t, v.Date.Time, time.Date(2019, 11, 01, 17, 22, 0, 0, v.Date.Location()))

	var categories []Category
	categories = append(categories, Category{
		Name: "Бытовая техника",
		ID:   "1",
	})
	categories = append(categories, Category{
		Name:     "Мелкая техника для кухни",
		ID:       "10",
		ParentId: "1",
	})
	categories = append(categories, Category{
		Name:     "Сэндвичницы и приборы для выпечки",
		ID:       "101",
		ParentId: "10",
	})
	assert.ElementsMatch(t, v.Shop.Categories, categories, "Categories not matched")

	var currencies []Currency
	currencies = append(currencies, Currency{
		ID:   "RUR",
		Rate: 1,
		Plus: 0,
	})
	currencies = append(currencies, Currency{
		ID:   "USD",
		Rate: 60,
		Plus: 0,
	})
	assert.ElementsMatch(t, v.Shop.Currencies, currencies, "Currencies not matched")

	var offers []Offer

	var deliveryOptions []Option
	deliveryOptions = append(deliveryOptions, Option{
		Cost: 300,
		Days: "1-3",
	})
	var pickupOptions []Option
	pickupOptions = append(pickupOptions, Option{
		Cost:        350,
		Days:        "1",
		OrderBefore: 12,
	})

	var params []Param
	params = append(params, Param{
		Name:  "Мощность",
		Value: "750 Вт",
	})
	offers = append(offers, Offer{
		Type:                 "vendor.model",
		Model:                "K220Y9",
		Vendor:               "Brand",
		TypePrefix:           "Сэндвичница",
		VendorCode:           "A1234567B",
		ID:                   "12346",
		BID:                  "60",
		Url:                  "http://best.seller.ru/product_page.asp?pid=12345",
		Price:                1099,
		OldPrice:             1399,
		EnableAutoDiscounts:  false,
		CurrencyId:           "RUR",
		CategoryId:           "101",
		Picture:              "http://best.seller.ru/img/device56789.jpg",
		Delivery:             true,
		Pickup:               true,
		DeliveryOptions:      deliveryOptions,
		PickupOptions:        pickupOptions,
		Description:          "Сэндвичница 2 в 1: можно приготовить как сэндвичи, так и вафли.",
		Store:                false,
		SalesNotes:           "Наличные, Visa/Mastercard, б/н расчет",
		ManufacturerWarranty: true,
		CountryOfOrigin:      "Россия",
		Barcode:              "9876543210",
		Param:                params,
		Weight:               "1.03",
		Dimensions:           "20.800/23.500/9.000",
	})
	assert.ElementsMatch(t, v.Shop.Offers, offers)
}
