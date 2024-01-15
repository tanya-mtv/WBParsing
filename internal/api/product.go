package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"parsingWB/internal/config"
	"parsingWB/internal/logger"
	"parsingWB/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ProductService struct {
	cfg        *config.Config
	log        logger.Logger
	httpClient *retryablehttp.Client
	repository storage
}

func NewProductService(cfg *config.Config, log logger.Logger, repo storage) *ProductService {

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = RetryMax
	retryClient.RetryWaitMin = RetryWaitMin
	retryClient.RetryWaitMax = RetryWaitMax
	retryClient.Backoff = backoff

	return &ProductService{
		cfg:        cfg,
		httpClient: retryClient,
		log:        log,
		repository: repo,
	}
}

func (p *ProductService) Post(ctx context.Context) (models.Out, error) {
	// var data models.Data
	data := models.Out{}
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", p.cfg.WbListUrl, bytes.NewBuffer([]byte(prepareAPIBody(p.cfg.Limit))))
	if err != nil {
		p.log.Errorf(err.Error())
		return data, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", p.cfg.Token)

	resp, err := p.httpClient.Do(req)

	if err != nil {
		p.log.Errorf(err.Error())
		return data, err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Errorf(err.Error())
		return data, err
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		p.log.Errorf(err.Error())

		return data, err
	}

	// fmt.Printf("Data %+v\n", data)
	return data, err
}

func (p *ProductService) PostPagination(ctx context.Context) ([]models.Out, error) {

	var list []models.Out
	var data models.Out
	for {
		req, err := retryablehttp.NewRequestWithContext(ctx, "POST", p.cfg.WbListUrl, bytes.NewBuffer([]byte(prepareAPIBody(p.cfg.Limit))))
		if err != nil {
			p.log.Errorf("Post erro to WB", err.Error())
			return list, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Authorization", p.cfg.Token)

		resp, err := p.httpClient.Do(req)

		if err != nil {
			p.log.Errorf("Can't do http request", err.Error())
			return list, err
		}
		defer resp.Body.Close()

		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			p.log.Errorf("Can't read response body", err.Error())
			return list, err
		}

		if err := json.Unmarshal(jsonData, &data); err != nil {
			p.log.Errorf("Can't unmarshal body", err.Error())
			fmt.Println("JSON body: ", string(jsonData))

			return list, err
		}

		colRows := data.Data.Cursor.Total
		if colRows <= p.cfg.Limit {
			list = append(list, data)
			break
		}

	}
	return list, nil

}

func (p *ProductService) ParsePage(ctx context.Context, card models.Cards) models.Product {
	var product models.Product
	product.Name = card.Title
	product.NmID = card.NmID

	// initialize a Chrome browser instance on port 4444
	service, err := selenium.NewChromeDriverService(p.cfg.ChromeDriver, p.cfg.ChromePort)
	if err != nil {
		p.log.Fatal("ERROR: ", err)
	}
	defer service.Stop()

	// configure the browser options
	// see
	// https://stackoverflow.com/questions/50642308/webdriverexception-unknown-error-devtoolsactiveport-file-doesnt-exist-while-t
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"start-maximized",         // open Browser in maximized mode
		"disable-infobars",        // disabling infobars
		"--headless",              // comment out this line for testing
		"--disable-extensions",    // disabling extensions
		"--disable-dev-shm-usage", // overcome limited resource problems
		"--disable-gpu",           // applicable to windows os only
		"--no-sandbox",            // Bypass OS security model
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		p.log.Errorf("ERROR: ", err)
		// return
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		p.log.Errorf("ERROR: ", err)
	}

	url := fmt.Sprintf("%s%d/detail.aspx", p.cfg.WbCatalogUrl, card.NmID)
	err = driver.Get(url)
	if err != nil {
		p.log.Errorf("ERROR: ", err)
	}

	time.Sleep(3 * time.Second)

	priceElements, err := driver.FindElements(selenium.ByCSSSelector, ".price-block__content")
	if err != nil {
		p.log.Errorf("ERROR: ", err)
	}
	for _, priceElements := range priceElements {
		priceElement, err := priceElements.FindElement(selenium.ByCSSSelector, "ins.price-block__final-price")
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		price, err := priceElement.Text()
		// fmt.Println("nameElement ", nameElement, "name ", name, "price ", price)

		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		if price == "" {
			continue
		}

		priceCnt := strings.Replace(strings.Replace(price, "₽", "", -1), " ", "", -1)
		priceFloat, err := strconv.ParseFloat(priceCnt, 64)

		if err != nil {
			fmt.Println("Error parsing price string to float", err.Error())
		}
		product.Price = priceFloat
	}

	productElements, err := driver.FindElements(selenium.ByCSSSelector, ".bestsellers__item-wrap")
	if err != nil {
		p.log.Errorf("Can't find bestsellers tag:", err)
		return product
	}

	// iterate over the product elements
	// and extract data from them
	if len(productElements) == 0 {
		p.log.Debug("Can't find another seler")
		return product
	}

	var sellerPrice []models.SellerPrice

	for _, productElement := range productElements {
		// select the name and price nodes
		nameElement, err := productElement.FindElement(selenium.ByCSSSelector, "p.bestsellers__name")
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}
		priceElement, err := productElement.FindElement(selenium.ByCSSSelector, "span.button-basket__btn-text")
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		// extract the data of interest
		name, err := nameElement.Text()
		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		price, err := priceElement.Text()

		if err != nil {
			fmt.Println("ERROR: ", err)
			continue
		}

		if price == "" || name == "" {
			continue
		}

		var sp models.SellerPrice

		priceCnt := strings.Replace(strings.Replace(price, "₽", "", -1), " ", "", -1)
		priceFloat, err := strconv.ParseFloat(priceCnt, 64)
		if err != nil {
			p.log.Errorf("Error parsing price string to float", err.Error())
			return product
		}

		sp = models.SellerPrice{
			Saller: name,
			Price:  priceFloat,
		}

		sellerPrice = append(sellerPrice, sp)

	}

	product.SellerPrice = sellerPrice

	err = p.repository.InsertData(product)
	if err != nil {
		p.log.Errorf("Can't insert data to DB")
	}
	return product

}

func backoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	sleepTime := min + min*time.Duration(2*attemptNum)
	return sleepTime
}
