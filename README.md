# Pepper Breeding Probability Calculator

A Go program that calculates trait inheritance probabilities in pepper breeding projects. This tool helps breeders predict the distribution of traits across generations and plan their breeding programs effectively.

## Features

- Calculates F2 generation probabilities for multiple traits
- Supports both dominant and recessive trait inheritance
- Provides probability ratios and expected plant numbers
- Sorts results by probability for easy interpretation
- Handles multiple plant traits simultaneously

## Genetic Background

The calculator works with three main inheritance patterns:
- Leaf Shape: Regular (L) dominant over Mutant (l)
- Foliage Color: Green (C) dominant over Purple (c)
- Fruit Shape: Long (F) dominant over Round (f)

### Inheritance Patterns

F1 Generation:
- Shows dominant traits from both parents
- All plants look identical
- Used to determine which traits are dominant

F2 Generation Distribution (64 plants):
- 27/64 (42.2%): All dominant traits
- 9/64 (14.1%): Single recessive trait
- 3/64 (4.7%): Double recessive traits
- 1/64 (1.6%): Triple recessive traits

## Usage

```go
// Define parent plants with their traits
purpleFlash := Plant{
    Name: "Purple Flash",
    Traits: []Trait{
        {Name: "Regular Leave Shape", Dominant: true, GeneLabel: "L"},
        {Name: "Purple Foliage", Dominant: false, GeneLabel: "C"},
        {Name: "Round Fruit Shape", Dominant: false, GeneLabel: "F"},
    },
}

candlelight := Plant{
    Name: "Candlelight Mutant",
    Traits: []Trait{
        {Name: "Mutant Leave Shape", Dominant: false, GeneLabel: "L"},
        {Name: "Green Foliage", Dominant: true, GeneLabel: "C"},
        {Name: "Long Fruit Shape", Dominant: true, GeneLabel: "F"},
    },
}

// Calculate probabilities for a specific number of plants
combinations := calculateF2Probabilities(purpleFlash, candlelight, 32)
```

## Example Output

```
F2 Generation Probabilities for Purple Flash × Candlelight Mutant
Total plants: 32
=====================================================
27/64 (42.2%) = Regular Leave Shape, Green Foliage, Long Fruit Shape
    Genotype: L_ C_ F_
    Expected number of plants: 13.5

9/64 (14.1%) = Mutant Leave Shape, Green Foliage, Long Fruit Shape
    Genotype: ll C_ F_
    Expected number of plants: 4.5

[... additional combinations ...]
```

## Structure

### Trait Type
```go
type Trait struct {
    Name      string  // Name of the trait
    Dominant  bool    // Whether the trait is dominant
    GeneLabel string  // Single letter used in genetic notation
}
```

### Plant Type
```go
type Plant struct {
    Name   string
    Traits []Trait
}
```

### TraitCombination Type
```go
type TraitCombination struct {
    Traits       []bool    // Trait expressions
    Probability  int       // Numerator of probability fraction
    Denominator  int       // Denominator (usually 64)
    Description  string    // Human-readable description
    GeneNotation string    // Genetic notation
    Expected     float64   // Expected number of plants
}
```

## Practical Applications

1. Breeding Program Planning:
   - Determine required population sizes
   - Estimate probability of desired combinations
   - Plan selection strategies

2. Resource Management:
   - Calculate space requirements
   - Plan growing areas
   - Optimize selection process

3. Trait Selection:
   - Identify rare combinations
   - Focus on desired trait combinations
   - Track inheritance patterns

## Limitations

- Assumes simple dominant/recessive inheritance
- Does not account for:
  - Incomplete dominance
  - Codominance
  - Multiple gene interactions
  - Environmental effects

## Future Improvements

- Add support for incomplete dominance
- Add visualization of results
- Support for larger trait sets
- Add F3 generation predictions
- Include confidence intervals
