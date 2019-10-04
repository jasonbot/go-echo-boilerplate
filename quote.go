package quoteapi

type quoteData struct {
	QuoteID  uint64 `json:"id"`
	Text     string `json:"text"`
	YodaText string `json:"yoda_text"`
}
