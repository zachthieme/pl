package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
  "time"
  "sync"
	"github.com/fogleman/gg"
)
// structs for import
type Order struct {
    ProducerOrderId       string           `json:"producerOrderId"`
    DeliveryShortOrderId  *string          `json:"deliveryShortOrderId"`
    ProviderOrderId       *string          `json:"providerOrderId"`
    ProviderId            string           `json:"providerId"`
    Channel               string           `json:"channel"`
    PickupIntent          string           `json:"pickupIntent"`
    LocationNumber        string           `json:"locationNumber"`
    MerchantId            string           `json:"merchantId"`
    Timestamp             time.Time        `json:"timestamp"`
    Customer              Customer         `json:"customer"`
    OrderPreparation      string           `json:"orderPreparation"`
    Items                 []Item           `json:"items"`
    Sequencing            Sequencing       `json:"sequencing"`
    DeliveryAttempts      int              `json:"deliveryAttempts"`
    Bill                  Bill             `json:"bill"`
    OrderNotes            *string          `json:"orderNotes"`
}

type Customer struct {
    CustomerXID       *string `json:"customerXID"`
    DisplayName       string  `json:"displayName"`
    LastInitial       string  `json:"lastInitial"`
    RewardTierNumber  *int    `json:"rewardTierNumber"`
}

type Item struct {
    ChildItems                []ChildItem `json:"childItems"`
    CustomInstructions        *string     `json:"customInstructions"`
    SumOfPricesAfterDiscounts int         `json:"sumOfPricesAfterDiscounts"`
    SkuNumber                 string      `json:"skuNumber"`
    Quantity                  int         `json:"quantity"`
    PriceAmount               int         `json:"priceAmount"`
    DiscountAmount            int         `json:"discountAmount"`
    PriceAfterDiscount        int         `json:"priceAfterDiscount"`
}

type ChildItem struct {
    SkuNumber          string `json:"skuNumber"`
    Quantity           int    `json:"quantity"`
    PriceAmount        *int   `json:"priceAmount"`
    DiscountAmount     int    `json:"discountAmount"`
    PriceAfterDiscount *int   `json:"priceAfterDiscount"`
}

type Sequencing struct {
    TravelTimeDetails []TravelTimeDetail `json:"travelTimeDetails"`
}

type TravelTimeDetail struct {
    Type           string    `json:"type"`
    MinTravelTime  int       `json:"minTravelTime"`
    MaxTravelTime  int       `json:"maxTravelTime"`
    TransportMode  string    `json:"transportMode"`
    Timestamp      time.Time `json:"timestamp"`
}

type Bill struct {
    ReceiptNumber    string     `json:"receiptNumber"`
    SubtotalAmount   int        `json:"subtotalAmount"`
    TotalTaxAmount   int        `json:"totalTaxAmount"`
    TotalAmount      int        `json:"totalAmount"`
    TaxLabels        []TaxLabel `json:"taxLabels"`
    Tenders          []Tender   `json:"tenders"`
}

type TaxLabel struct {
    TaxLabel   string `json:"taxLabel"`
    TaxAmount  int    `json:"taxAmount"`
}

type Tender struct {
    TenderType   string `json:"tenderType"`
    Amount       int    `json:"amount"`
    CurrencyCode string `json:"currencyCode"`
}

func main() {

	// Read JSON from args
  jsonFile, err := os.Open(os.Args[1])
  if err != nil {
    panic(err)
  }
  defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
    log.Fatalf("Error reading from stdin: %v", err)
	}

	var order Order
	if err := json.Unmarshal(bytes, &order); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

  var starttime = time.Now().UnixNano() / int64(time.Millisecond)

	var wg sync.WaitGroup

	// Generate labels for each item using goroutines
	for i, item := range order.Items {
		wg.Add(1)
		go createLabel(i, order.Customer, item, order, &wg)
	}

	// Wait for all goroutines to complete
	wg.Wait()
  var endtime = time.Now().UnixNano() / int64(time.Millisecond)
  var diff = endtime - starttime 
  fmt.Printf("All labels created: time (ms): %d\n", diff)
}

func createLabel(index int, customer Customer, item Item, order Order, wg *sync.WaitGroup) {
	defer wg.Done()  
	const width = 150
	var height = 150 + 40*len(item.ChildItems)
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	fontFace,err := gg.LoadFontFace("FreeMono.ttf", 7)
  if err != nil {
		log.Fatalf("Error loading font: %v", err)
	}

  y := 30

  dc.DrawString(fmt.Sprintf("Item: %d of %d ", index + 1, len(order.Items)), 20, float64(y))
  y += 20

  dc.DrawString(fmt.Sprintf("Items in order: %d ", len(order.Items)), 20, float64(y))
  y += 40

	// Write customer info using different font
 	fontFace, _ = gg.LoadFontFace("DejaVuSerif.ttf", 24 )
  if err != nil {
		log.Fatalf("Error loading font: %v", err)
	}

	dc.SetFontFace(fontFace)
	dc.DrawString(fmt.Sprintf("%s %s.", customer.DisplayName, customer.LastInitial), 20, float64(y))
  y += 30
    
  // Write item info - at least what i had in the JSON
 	fontFace, _ = gg.LoadFontFace("FreeSans.ttf", 12 )
  if err != nil {
		log.Fatalf("Error loading font: %v", err)
	}

	dc.SetFontFace(fontFace)
	dc.DrawString(fmt.Sprintf("%s", item.SkuNumber), 20,  float64(y))
	y += 20

	// Write child items info
	for _, child := range item.ChildItems {
		dc.DrawString(fmt.Sprintf("%s", child.SkuNumber), 30, float64(y))
		y += 20
	}
	
	fileName := fmt.Sprintf("label_%d.png", index)
	dc.SavePNG(fileName)
  fmt.Printf("Label saved as: %s\n ", fileName)
 }
