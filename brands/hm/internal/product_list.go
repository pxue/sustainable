package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ProductListService service

type FilterValue struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Filtercount int    `json:"filtercount"`
	Code        string `json:"code"`
	Selected    bool   `json:"selected"`
	Disabled    bool   `json:"disabled"`
}

type FilterGroup struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	FilterValues []*FilterValue `json:"filtervalues"`
}

type Filter struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	FilterValues []*FilterValue `json:"filtervalues,omitempty"`
	Group        []*FilterGroup `json:"group,omitempty"`
}

type Sort struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Filtervalues []*FilterValue `json:"filtervalues"`
}

type ProductList struct {
	Total      int                `json:"total"`
	ItemsShown int                `json:"itemsShown"`
	Filters    []*Filter          `json:"filters"`
	SortBy     *Sort              `json:"sortby"`
	Products   []*ProductListItem `json:"products"`
	Labels     struct {
		FilterBy      string `json:"filterBy"`
		TotalCount    string `json:"totalCount"`
		ShowItemsText string `json:"showItemsText"`
		LoadMoreText  string `json:"loadMoreText"`
	} `json:"labels"`
	Datatracking struct {
		FilterUsed       string `json:"filterUsed"`
		FilterChanged    string `json:"filterChanged"`
		FilterRemoved    string `json:"filterRemoved"`
		LoadMoreProducts string `json:"loadMoreProducts"`
	} `json:"datatracking"`
}

type ProductListItem struct {
	ArticleCode               string              `json:"articleCode"`
	OnClick                   string              `json:"onClick"`
	Link                      string              `json:"link"`
	Title                     string              `json:"title"`
	Category                  string              `json:"category"`
	Image                     []*ProductListImage `json:"image"`
	LegalText                 string              `json:"legalText"`
	PromotionalMarkerText     string              `json:"promotionalMarkerText"`
	ShowPromotionalClubMarker bool                `json:"showPromotionalClubMarker"`
	ShowPriceMarker           bool                `json:"showPriceMarker"`
	FavouritesTracking        string              `json:"favouritesTracking"`
	FavouritesSavedText       string              `json:"favouritesSavedText"`
	FavouritesNotSavedText    string              `json:"favouritesNotSavedText"`
	MarketingMarkerText       string              `json:"marketingMarkerText"`
	MarketingMarkerType       string              `json:"marketingMarkerType"`
	MarketingMarkerCSS        string              `json:"marketingMarkerCss"`
	Price                     string              `json:"price"`
	RedPrice                  string              `json:"redPrice"`
	YellowPrice               string              `json:"yellowPrice"`
	BluePrice                 string              `json:"bluePrice"`
	ClubPriceText             string              `json:"clubPriceText"`
	SellingAttribute          string              `json:"sellingAttribute"`
	SwatchesTotal             string              `json:"swatchesTotal"`
	Swatches                  []struct {
		ColorCode   string `json:"colorCode"`
		ArticleLink string `json:"articleLink"`
		ColorName   string `json:"colorName"`
	} `json:"swatches"`
	OutOfStockText string `json:"outOfStockText"`
	ComingSoon     string `json:"comingSoon"`
}

type ProductListImage struct {
	Src          string `json:"src"`
	DataAltImage string `json:"dataAltImage"`
	Alt          string `json:"alt"`
	DataAltText  string `json:"dataAltText"`
}

type ProductListParams struct {
	Sort      string
	ImageSize string
	Image     string
	Offset    int
	PageSize  int
}

func (p *ProductListParams) Encode() string {
	value := url.Values{}

	value.Add("sort", p.Sort)
	value.Add("image-size", p.ImageSize)
	value.Add("image", p.Image)
	value.Add("offset", fmt.Sprintf("%d", p.Offset))
	value.Add("page-size", fmt.Sprintf("%d", p.PageSize))

	return value.Encode()
}

func (s *ProductListService) ShopByProduct(ctx context.Context, params *ProductListParams) ([]*ProductListItem, *http.Response, error) {
	path, _ := url.Parse(fmt.Sprintf("/%s/men/shop-by-product/view-all/_jcr_content/main/productlisting_fa5b.display.json", s.client.opts.countryCode))
	baseUrl, _ := url.Parse("https://www2.hm.com")
	u := baseUrl.ResolveReference(path)
	u.RawQuery = params.Encode()

	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var wrapper *ProductList
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper.Products, resp, nil
}
