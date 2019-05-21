package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type DraftCard struct {
	Name  string
	Value float64
}

type EternalCard struct {
	SetNumber     int     `json:"SetNumber"`
	EternalId     int     `json:"EternalID"`
	Name          string  `json:"Name"`
	Influence     string  `json:"Influence"`
	Type          string  `json:"Type"`
	LimitedValue  float64 `json:"LimitedValue"`
	Rarity        string
	DeckBuildable bool
}

func main() {
	// our running set of card data
	data := make(map[string]DraftCard)

	records := getData("data/Set6DraftTierSummary.csv")
	for _, r := range records {
		v, _ := strconv.ParseFloat(r[1], 64)
		data[r[0]] = DraftCard{r[0], v}
	}

	records = getData("data/Set5.5DraftTierSummary.csv")
	for _, r := range records {
		if _, ok := data[r[0]]; !ok {
			v, _ := strconv.ParseFloat(r[1], 64)
			data[r[0]] = DraftCard{r[0], v}
		}
	}

	records = getData("data/Set5DraftTierSummary.csv")
	for _, r := range records {
		if _, ok := data[r[0]]; !ok {
			v, _ := strconv.ParseFloat(r[1], 64)
			data[r[0]] = DraftCard{r[0], v}
		}
	}

	records = getData("data/Set4DraftTierSummary.csv")
	for _, r := range records {
		if _, ok := data[r[0]]; !ok {
			v, _ := strconv.ParseFloat(r[1], 64)
			data[r[0]] = DraftCard{r[0], v}
		}
	}

	records = getData("data/Set3DraftTierSummary.csv")
	for _, r := range records {
		if _, ok := data[r[0]]; !ok {
			v, _ := strconv.ParseFloat(r[1], 64)
			data[r[0]] = DraftCard{r[0], v}
		}
	}

	records = getData("data/InfluenceGivers.csv")
	for _, r := range records {
		if _, ok := data[r[0]]; !ok {
			v, _ := strconv.ParseFloat(r[1], 64)
			data[r[0]] = DraftCard{r[0], v}
		}
	}

	bytes, _ := ioutil.ReadFile("data/eternal-cards.json")
	var cards []EternalCard
	err := json.Unmarshal(bytes, &cards)
	if err != nil {
		println(err.Error())
	}

	var notFoundCards []EternalCard
	b := cards[:0]
	for _, c := range cards {
		if _, ok := data[c.Name]; ok {
			c.LimitedValue = data[c.Name].Value
			c.Influence = strings.Replace(c.Influence, "{", "", -1)
			c.Influence = strings.Replace(c.Influence, "}", "", -1)
			b = append(b, c)
		} else if c.SetNumber > 0 && c.SetNumber < 1000 &&
			c.Rarity != "Promo" && c.Rarity != "None" {
			// Set 0 isn't in packs
			// Set 1000+ are campaigns
			// Promos aren't in packs
			// Rarity == None denotes generated cards
			notFoundCards = append(notFoundCards, c)
		}
	}

	json, err := json.MarshalIndent(b, "", "  ")
	fmt.Println(string(json))

	// uncomment to debug if any cards are not being found
	if len(notFoundCards) > 0 {
		fmt.Printf("%+v\n", notFoundCards)
	}
}

func getData(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		println(err.Error())
	}
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		println(err.Error())
	}
	return records
}
