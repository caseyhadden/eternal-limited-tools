package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type EternalCard struct {
	SetNumber    int    `json:"SetNumber"`
	EternalId    int    `json:"EternalID"`
	Name         string `json:"Name"`
	CardText     string `json:"CardText"`
	Cost         int    `json:"Cost"`
	Influence    string `json:"Influence"`
	Attack       int    `json:"Attack"`
	Health       int    `json:"Health"`
	Rarity       string `json:"Rarity"`
	Type         string `json:"Type"`
	LimitedValue float64
}

type ExportedCard struct {
	Name       string
	SetNumber  int
	CardNumber int
}

var exportedPool []ExportedCard
var allCards []EternalCard

func main() {
	fileName := "pool.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	f, _ := os.Open(fileName)
	defer f.Close()

	var exportedPool []ExportedCard
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		card := ExportedCard{}
		// like - 1 Deft Strike (Set4 #71)
		line := scanner.Text()
		values := strings.Split(line, "(")
		numberOfAndName := values[0]
		setAndNumber := values[1]
		numberOfCards, _ := strconv.Atoi(numberOfAndName[:1])
		card.Name = strings.TrimSpace(numberOfAndName[2:])
		fmt.Sscanf(setAndNumber, "Set%d #%d)", &card.SetNumber, &card.CardNumber)
		for i := 0; i < numberOfCards; i++ {
			exportedPool = append(exportedPool, card)
		}
	}

	bytes, _ := ioutil.ReadFile("cards-and-values.json")
	json.Unmarshal(bytes, &allCards)

	var poolCards []EternalCard
	for _, v := range exportedPool {
		c := findCard(v)
		if c.EternalId != 0 {
			poolCards = append(poolCards, findCard(v))
		}
	}

	sort.Slice(poolCards, func(i, j int) bool {
		return poolCards[i].LimitedValue < poolCards[j].LimitedValue
	})

	var total float64 = 0
	for _, v := range poolCards {
		fmt.Println(fmt.Sprintf("%s - %s - %f", v.Name, v.Influence, v.LimitedValue))
		total += v.LimitedValue
	}
	fmt.Println(fmt.Sprintf("Average - %f", total/float64(len(poolCards))))
}

func findCard(c ExportedCard) EternalCard {
	for _, v := range allCards {
		if c.SetNumber == v.SetNumber &&
			c.CardNumber == v.EternalId {
			return v
		}
	}
	fmt.Printf("%+v\n", c)
	return EternalCard{}
}
