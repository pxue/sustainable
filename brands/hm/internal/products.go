package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ProductService service

type Product struct {
	Code                string        `json:"code"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	SapProductName      string        `json:"sapProductName"`
	SellingAttributes   []string      `json:"sellingAttributes"`
	Color               *Color        `json:"color"`
	WhitePrice          *Price        `json:"whitePrice"`
	PriceType           string        `json:"priceType"`
	ImportedBy          string        `json:"importedBy"`
	ImportedDate        string        `json:"importedDate"`
	NetQuantity         string        `json:"netQuantity"`
	CountryOfProduction string        `json:"countryOfProduction"`
	ProductTypeName     string        `json:"productTypeName"`
	Measurements        []interface{} `json:"measurements"`
	DescriptiveLenght   []interface{} `json:"descriptiveLenght"`
	Fits                []string      `json:"fits"`
	ShowPriceMarker     bool          `json:"showPriceMarker"`
	BaseProductCode     string        `json:"baseProductCode"`
	AncestorProductCode string        `json:"ancestorProductCode"`
	MainCategory        struct {
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"mainCategory"`
	Supercategories []struct {
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"supercategories"`
	ConstructionDescr string            `json:"constructionDescr"`
	CustomerGroup     string            `json:"customerGroup"`
	Functions         []interface{}     `json:"functions"`
	NewArrival        bool              `json:"newArrival"`
	ArticlesList      []*Variant        `json:"articlesList"`
	InStock           bool              `json:"inStock"`
	ProductURL        string            `json:"productUrl"`
	SwatchesType      string            `json:"swatchesType"`
	RootCategoryPath  string            `json:"rootCategoryPath"`
	Styles            []interface{}     `json:"styles"`
	MaterialDetails   []*MaterialDetail `json:"materialDetails"`
	PresentationTypes string            `json:"presentationTypes"`
	NewProduct        bool              `json:"newProduct"`
}

type Color struct {
	Code     string `json:"code"`
	Text     string `json:"text"`
	RgbColor string `json:"rgbColor"`
}

type Price struct {
	Price         float64 `json:"price"`
	Currency      string  `json:"currency"`
	ReferenceFlag bool    `json:"referenceFlag"`
	StartDate     int64   `json:"startDate"`
	EndDate       int64   `json:"endDate"`
}

type MaterialDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StyleWith struct {
	Code                       string `json:"code"`
	ShowPriceMarker            bool   `json:"showPriceMarker"`
	StyleWithOrigin            string `json:"styleWithOrigin"`
	ParentProductCode          string `json:"parentProductCode"`
	InStore                    bool   `json:"inStore"`
	ProductTransparencyEnabled bool   `json:"productTransparencyEnabled"`
	SuppliersDetailEnabled     bool   `json:"suppliersDetailEnabled"`
}

type Asset struct {
	URL       string `json:"url"`
	AssetType string `json:"assetType"`
}

type Variant struct {
	*Product

	ColourDescription              string         `json:"colourDescription"`
	Pattern                        string         `json:"pattern"`
	GalleryDetails                 []*Asset       `json:"galleryDetails"`
	FabricSwatchThumbnails         []*Asset       `json:"fabricSwatchThumbnails"`
	StyleWith                      []*StyleWith   `json:"styleWith,omitempty"`
	CareInstructions               []string       `json:"careInstructions"`
	Compositions                   []*Composition `json:"compositions"`
	GraphicalAppearanceDesc        string         `json:"graphicalAppearanceDesc"`
	GenericDescription             string         `json:"genericDescription"`
	VariantsList                   []*SizeVariant `json:"variantsList,omitempty"`
	Concepts                       []string       `json:"concepts"`
	LegalRestrictions              []interface{}  `json:"legalRestrictions"`
	ParentProductCode              string         `json:"parentProductCode"`
	StyleWithScenario              string         `json:"styleWithScenario"`
	InStore                        bool           `json:"inStore"`
	ProductTransparencyEnabled     bool           `json:"productTransparencyEnabled"`
	SuppliersDetailEnabled         bool           `json:"suppliersDetailEnabled"`
	SuppliersSectionDisabledReason string         `json:"suppliersSectionDisabledReason"`
}

type Composition struct {
	Materials []struct {
		Name       string `json:"name"`
		Percentage string `json:"percentage"`
	} `json:"materials"`
	CompositionType string `json:"compositionType,omitempty"`
}

type Size struct {
	SizeCode      string `json:"sizeCode"`
	Name          string `json:"name"`
	SizeScaleCode string `json:"sizeScaleCode"`
	SizeOrder     int    `json:"sizeOrder"`
	SizeFilter    string `json:"sizeFilter"`
	Market        string `json:"market"`
}

type SizeVariant struct {
	Code   string `json:"code"`
	Size   *Size  `json:"size,omitempty"`
	Width  string `json:"width,omitempty"`
	Length string `json:"length,omitempty"`
}

func (s *ProductService) Get(ctx context.Context, productID string) (*Product, *http.Response, error) {
	path, _ := url.Parse(fmt.Sprintf("/hmwebservices/service/article/get-article-by-code/%s/Online/%s/en.json", s.client.opts.country, productID))
	baseUrl, _ := url.Parse("https://app2.hm.com")
	u := baseUrl.ResolveReference(path)

	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	wrapper := struct {
		Product *Product `json:"product"`
	}{}
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Product, resp, nil
}
