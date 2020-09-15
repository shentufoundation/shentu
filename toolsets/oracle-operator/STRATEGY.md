# Aggregation Strategies

Aggregation method for primitive scores.

## `linear`

Linear combination (with weights) of primitive scores.

```toml
[strategy.eth]
# combination strategy
type = "linear"
[[strategy.eth.primitive]]
primitive_contract_address = "certik1r4834vyyu8vrarxgyatn34j8lsguyhn7csl0ju"
weight = 0.1
[[strategy.eth.primitive]]
primitive_contract_address = "certik1r4834vyyu8vrarxgyatn34j8lsguyhn7csl0ju"
weight = 0.1
```
