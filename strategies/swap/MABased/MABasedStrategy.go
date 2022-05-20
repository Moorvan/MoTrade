package MABased

import "MoTrade/strategies"

type MABasedStrategy struct {
	strategies.Strategy
}

func NewMABasedStrategy() *MABasedStrategy {
	return &MABasedStrategy{}
}
