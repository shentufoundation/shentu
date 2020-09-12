package oracle

import (
	"fmt"
	"strings"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

// Aggregation for scores.
type Aggregation interface {
	Aggregate(scores <-chan types.PrimitiveScore) (uint8, error)
}

// NewAggregation returns the strategy implementation based on configured type.
func NewAggregation(strategy types.Strategy) (Aggregation, error) {
	switch strings.ToLower(strategy.Type) {
	case "linear":
		return LinearCombination{}, nil
	default:
		return nil, fmt.Errorf("Aggregation type %s not found", strategy.Type)
	}
}

// LinearCombination of scores.
type LinearCombination struct {
	Aggregation
}

// Aggregate performs linear combination of the primitive scores.
func (LinearCombination) Aggregate(scores <-chan types.PrimitiveScore) (uint8, error) {
	if len(scores) == 0 {
		return 0, fmt.Errorf("no primive score")
	}
	sum := float32(0)
	weightSum := float32(0)
	for score := range scores {
		sum += float32(score.Score) * score.Primitive.Weight
		weightSum += score.Primitive.Weight
	}
	if weightSum == 0 {
		return 0, fmt.Errorf("linear combination weight error")
	}
	return uint8(sum / weightSum), nil
}
