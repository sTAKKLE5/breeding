# Pepper Breeding Probability Calculator

A Go program that calculates trait inheritance probabilities in pepper breeding projects. This tool helps breeders predict the distribution of traits across generations and plan their breeding programs effectively.

## Features

- Calculates F2 generation probabilities for multiple traits
- Supports both dominant and recessive trait inheritance
- Provides probability ratios and expected plant numbers
- Sorts results by probability for easy interpretation
- Handles multiple plant traits simultaneously
- Filters combinations based on target genotypes
- Calculates total probability for desired trait combinations
- Provides human-readable trait descriptions
- Supports filtering for both dominant and recessive traits

## Genetic Background

The calculator works with Mendelian inheritance patterns for multiple traits. Each trait is defined by:
- A dominant allele (typically represented by a capital letter, e.g., 'A')
- A recessive allele (typically represented by a lowercase letter, e.g., 'a')
- Whether the trait shows dominant or recessive inheritance

Common pepper traits that follow Mendelian inheritance include:
- Leaf shape variations
- Foliage coloration
- Fruit shape
- Plant height
- Fruit position
- And others

### Inheritance Patterns

F1 Generation:
- First generation after crossing two parent lines
- All plants show uniform appearance
- Reveals which traits are dominant/recessive
- Heterozygous for all differing parental traits

F2 Generation Distribution:
For each individual trait:
- 3/4 (75%) show dominant phenotype
- 1/4 (25%) show recessive phenotype

For multiple traits, the probabilities follow multiplicative inheritance:
- For n traits, there are 2^n possible combinations
- Probability of all dominant traits: (3/4)^n
- Probability of all recessive traits: (1/4)^n
- Other combinations follow binomial distribution

Example probabilities for different numbers of traits:
```
Single trait (n=1):
- Dominant (A_): 3/4
- Recessive (aa): 1/4

Two traits (n=2):
- Both dominant (A_B_): 9/16
- One recessive: 3/16 each
- Both recessive (aabb): 1/16

Three traits (n=3):
- All dominant (A_B_C_): 27/64
- Two dominant, one recessive: 9/64 each
- One dominant, two recessive: 3/64 each
- All recessive (aabbcc): 1/64
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

### TargetGenotype Type
```go
type TargetGenotype struct {
    Genotype    string    // Genetic notation (e.g., "ll" or "L")
    Description string    // Human-readable description
}
```

### SummaryStats Type
```go
type SummaryStats struct {
    TotalProbabilityNum   int     // Numerator of total probability
    TotalProbabilityDenom int     // Denominator of total probability
    Percentage            float64 // Percentage of total
    ExpectedPlants        float64 // Expected number of plants
}
```

## Usage

### Basic Usage
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

// Calculate all combinations
allCombinations := calculateF2Probabilities(purpleFlash, candlelight, 64)
```

### Filtering for Specific Traits
You can filter for specific trait combinations using either dominant or recessive notation:

```go
// For recessive traits (use doubled letters)
targetGenotypes := []TargetGenotype{
    {Genotype: "ll", Description: getGenotypeDescription("ll", purpleFlash, candlelight)}, // Mutant leaves
    {Genotype: "cc", Description: getGenotypeDescription("cc", purpleFlash, candlelight)}, // Purple foliage
}

// For dominant traits (use single letters)
targetGenotypes := []TargetGenotype{
    {Genotype: "L", Description: getGenotypeDescription("L", purpleFlash, candlelight)}, // Regular leaves
    {Genotype: "F", Description: getGenotypeDescription("F", purpleFlash, candlelight)}, // Long fruit
}

// Mix of dominant and recessive
targetGenotypes := []TargetGenotype{
    {Genotype: "L", Description: getGenotypeDescription("L", purpleFlash, candlelight)},  // Regular leaves
    {Genotype: "cc", Description: getGenotypeDescription("cc", purpleFlash, candlelight)}, // Purple foliage
}

// Filter combinations based on target genotypes
filteredCombinations, summary := filterCombinations(allCombinations, targetGenotypes)
```

## Genotype Notation

The calculator supports two types of genotype notation for filtering:

1. Recessive traits: Use doubled lowercase letters
   - `ll` for mutant leaves
   - `cc` for purple foliage
   - `ff` for round fruit

2. Dominant traits: Use single uppercase letters
   - `L` for regular leaves
   - `C` for green foliage
   - `F` for long fruit

The genotype notation in results uses:
- `L_` indicates presence of dominant allele (LL or Ll)
- `ll` indicates homozygous recessive
- Space separated for multiple traits (e.g., `L_ cc F_`)

## Example Outputs

### Filtering for Mixed Traits
```
F2 Generation Probabilities for Purple Flash × Candlelight Mutant
Total plants: 64

Target traits:
- Regular Leave Shape (L)
- Purple Foliage (cc)

Matching Combinations:
=====================
9/64 (14.1%) = Regular Leave Shape, Purple Foliage, Long Fruit Shape
    Genotype: L_ cc F_
    Expected number of plants: 9.0

3/64 (4.7%) = Regular Leave Shape, Purple Foliage, Round Fruit Shape
    Genotype: L_ cc ff
    Expected number of plants: 3.0

Summary Statistics:
==================
Total Probability: 12/64
Percentage: 18.8%
Expected Total Plants with Target Traits: 12.0
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

4. Target Trait Analysis:
   - Calculate probabilities for specific trait combinations
   - Determine required population sizes for rare combinations
   - Plan selective breeding programs
   - Estimate success rates for desired traits

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
- Add support for more complex filtering patterns
- Add visual representation of trait distributions
- Support for inheritance pattern validation
- Add batch processing capabilities

## License

MIT License
