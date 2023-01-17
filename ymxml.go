package YandexMarketXml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
)

type YandexMarketYML struct {
	FileName string
}

type YmlCatalog struct {
	YmlCatalog xml.Name   `xml:"yml_catalog"`
	Date       customTime `xml:"date,attr"`
	Shop       Shop       `xml:"shop"`
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02 15:04"
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return errors.Wrap(err, "decode element")
	}
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return errors.Wrap(err, "parse time")
	}
	*c = customTime{parse}
	return nil
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse("2006-01-02 15:04", attr.Value)
	if err != nil {
		return errors.Wrap(err, "parse time")
	}
	*c = customTime{parse}
	return nil
}

type Shop struct {
	Shop                xml.Name   `xml:"shop"`
	Name                string     `xml:"name"`
	Company             string     `xml:"company"`
	Url                 string     `xml:"url"`
	Platform            string     `xml:"platform,omitempty"`
	Version             string     `xml:"version,omitempty"`
	Agency              string     `xml:"agency,omitempty"`
	Email               string     `xml:"email,omitempty"`
	Currencies          []Currency `xml:"currencies>currency"`
	Categories          []Category `xml:"categories>category"`
	DeliveryOptions     []Option   `xml:"delivery-options"`
	PickupOptions       []Option   `xml:"pickup-options,omitempty"`
	EnableAutoDiscounts bool       `xml:"enable_auto_discounts,omitempty"`
	Offers              []Offer    `xml:"offers>offer"`
	Gifts               string     `xml:"gifts"`
	Promos              string     `xml:"promos"`
}

type Currency struct {
	ID   string  `xml:"id,attr"`
	Rate float64 `xml:"rate,attr"`
	Plus int64   `xml:"plus,attr"`
}

type Category struct {
	Name     string `xml:",chardata"`
	ID       string `xml:"id,attr"`
	ParentId string `xml:"parentId,attr,omitempty"`
}

type Option struct {
	Cost        int64  `xml:"cost,attr"`
	Days        string `xml:"days,attr"`
	OrderBefore int64  `xml:"order-before,attr,omitempty"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Unit  string `xml:"unit,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Offer struct {
	//Элементы, специфичные для произвольного типа описания
	Type       string `xml:"type,attr,omitempty"`
	Model      string `xml:"model,omitempty"`
	Vendor     string `xml:"vendor,omitempty"`
	TypePrefix string `xml:"typePrefix,omitempty"`
	//Стандартные элементы
	VendorCode           string     `xml:"vendorCode,omitempty"`
	ID                   string     `xml:"id,attr"`
	BID                  string     `xml:"bid,attr,omitempty"`
	Url                  string     `xml:"url"`
	Price                int64      `xml:"price"`
	OldPrice             int64      `xml:"oldprice,omitempty"`
	EnableAutoDiscounts  bool       `xml:"enable_auto_discounts,omitempty"`
	CurrencyId           string     `xml:"currencyId"`
	CategoryId           string     `xml:"categoryId"`
	Picture              string     `xml:"picture"`
	Delivery             bool       `xml:"delivery,omitempty"`
	Pickup               bool       `xml:"pickup,omitempty"`
	DeliveryOptions      []Option   `xml:"delivery-options>option"`
	PickupOptions        []Option   `xml:"pickup-options>option"`
	Store                bool       `xml:"store,omitempty"`
	Description          string     `xml:"description,omitempty"`
	SalesNotes           string     `xml:"sales_notes,omitempty"`
	MinQuantity          int64      `xml:"min-quantity,omitempty"`
	ManufacturerWarranty bool       `xml:"manufacturer_warranty,omitempty"`
	CountryOfOrigin      string     `xml:"country_of_origin,omitempty"`
	Adult                bool       `xml:"adult,omitempty"`
	Barcode              string     `xml:"barcode,omitempty"`
	Param                []Param    `xml:"param,omitempty"`
	Condition            string     `xml:"condition,omitempty"`       //TODO
	CreditTemplate       string     `xml:"credit-template,omitempty"` //TODO
	Expiry               customTime `xml:"expiry,omitempty"`
	Weight               string     `xml:"weight,omitempty"`
	Dimensions           string     `xml:"dimensions,omitempty"`
	Downloadable         bool       `xml:"downloadable,omitempty"`
	Available            bool       `xml:"available,attr,omitempty"`
	Age                  string     `xml:"age,omitempty"` //TODO
	GroupID              string     `xml:"group_id,attr,omitempty"`
}

func YandexMarketXML(fileName string) *YandexMarketYML {
	return &YandexMarketYML{FileName: fileName}
}

func (ymx *YandexMarketYML) CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func (ymx *YandexMarketYML) Parse() (YmlCatalog, error) {
	var ymlc YmlCatalog
	xmlFile, err := os.Open(ymx.FileName)
	if err != nil {
		return ymlc, errors.Wrap(err, "open file")
	}
	d := xml.NewDecoder(xmlFile)
	d.CharsetReader = ymx.CharsetReader
	err = d.Decode(&ymlc)
	if err != nil {
		return ymlc, errors.Wrap(err, "decode yaml")
	}
	return ymlc, err
}
