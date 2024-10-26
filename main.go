package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Trait struct {
	Name      string
	Dominant  bool
	GeneLabel string
}

type Plant struct {
	Name   string
	Traits []Trait
}

type TraitCombination struct {
	Traits       []bool
	Probability  int
	Denominator  int
	Description  string
	GeneNotation string
	Expected     float64
}

func calculateDenominator(numTraits int) int {
	return int(math.Pow(2, float64(2*numTraits)))
}

func calculateProbability(numTraits, numRecessive int) int {
	dominantTraits := numTraits - numRecessive
	return int(math.Pow(3, float64(dominantTraits)))
}

func calculateF2Probabilities(plant1, plant2 Plant, totalPlants int) []TraitCombination {
	numTraits := len(plant1.Traits)
	if numTraits != len(plant2.Traits) {
		panic("Plants must have same number of traits")
	}

	denominator := calculateDenominator(numTraits)
	numCombinations := int(math.Pow(2, float64(numTraits)))
	combinations := make([]TraitCombination, 0)

	for i := 0; i < numCombinations; i++ {
		traits := make([]bool, numTraits)
		description := make([]string, 0)
		geneNotation := make([]string, 0)
		numRecessive := 0

		for j := 0; j < numTraits; j++ {
			traits[j] = (i & (1 << j)) == 0
			if !traits[j] {
				numRecessive++
			}

			var traitName string
			if traits[j] {
				if plant1.Traits[j].Dominant {
					traitName = plant1.Traits[j].Name
				} else {
					traitName = plant2.Traits[j].Name
				}
				geneNotation = append(geneNotation, plant1.Traits[j].GeneLabel+"_")
			} else {
				if plant1.Traits[j].Dominant {
					traitName = plant2.Traits[j].Name
				} else {
					traitName = plant1.Traits[j].Name
				}
				geneNotation = append(geneNotation, strings.ToLower(plant1.Traits[j].GeneLabel+plant1.Traits[j].GeneLabel))
			}
			description = append(description, traitName)
		}

		probability := calculateProbability(numTraits, numRecessive)
		expected := float64(probability) * float64(totalPlants) / float64(denominator)

		combinations = append(combinations, TraitCombination{
			Traits:       traits,
			Probability:  probability,
			Denominator:  denominator,
			Description:  strings.Join(description, ", "),
			GeneNotation: strings.Join(geneNotation, " "),
			Expected:     expected,
		})
	}

	// Sort combinations by probability in descending order
	sort.Slice(combinations, func(i, j int) bool {
		return combinations[i].Probability > combinations[j].Probability
	})

	return combinations
}

func main() {
	purpleFlash := Plant{
		Name: "Purple Flash",
		Traits: []Trait{
			{
				Name:      "Regular Leave Shape",
				Dominant:  true,
				GeneLabel: "L",
			},
			{
				Name:      "Purple Flowers",
				Dominant:  false,
				GeneLabel: "A",
			},
			{
				Name:      "Purple Foliage",
				Dominant:  false,
				GeneLabel: "C",
			},
			{
				Name:      "Round Fruit Shape",
				Dominant:  false,
				GeneLabel: "F",
			},
		},
	}

	candlelight := Plant{
		Name: "Candlelight Mutant",
		Traits: []Trait{
			{
				Name:      "Mutant Leave Shape",
				Dominant:  false,
				GeneLabel: "L",
			},
			{
				Name:      "White Flowers",
				Dominant:  true,
				GeneLabel: "A",
			},
			{
				Name:      "Green Foliage",
				Dominant:  true,
				GeneLabel: "C",
			},
			{
				Name:      "Long Fruit Shape",
				Dominant:  true,
				GeneLabel: "F",
			},
		},
	}

	// Calculate for 32 plants
	totalPlants := 64

	combinations := calculateF2Probabilities(purpleFlash, candlelight, totalPlants)

	fmt.Printf("\nF2 Generation Probabilities for %s × %s\n", purpleFlash.Name, candlelight.Name)
	fmt.Printf("Total plants: %d\n\n", totalPlants)

	for _, combo := range combinations {
		percentage := float64(combo.Probability) / float64(combo.Denominator) * 100
		fmt.Printf("%d/%d (%0.1f%%) = %s\n",
			combo.Probability,
			combo.Denominator,
			percentage,
			combo.Description)
		fmt.Printf("    Genotype: %s\n", combo.GeneNotation)
		fmt.Printf("    Expected number of plants: %.1f\n\n", combo.Expected)
	}
}
