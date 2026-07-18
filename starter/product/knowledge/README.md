# Product Knowledge

Store durable knowledge owned by this product, including business rules, conventions, decisions, vocabulary, patterns, examples, and imported source evidence.

## Boundaries

- Product and cross-cutting decisions belong in `decisions/`; design decisions belong in `../design/decisions/`; engineering decisions belong in `../engineering/decisions/`. Every record is indexed in `.product/decisions.json` with its `domain` and `path`.
- Product gate commands and security conventions belong in `conventions/`.
- Existing source material and reviewed mappings belong in `imports/`.
- `assessments/product-landscape.md` is the complete code and product-boundary inventory. Keep it explicit about observed evidence versus hypotheses; finish it before delivery-domain modeling.
- Create other knowledge directories only when a real artifact needs them.

Framework maintenance history remains in the framework repository's Git history.
