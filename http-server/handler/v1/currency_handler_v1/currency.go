package currency_handler_v1

import (
	"net/http"
	"strconv"
	"strings"

	"paylist.server/middleware"
	"paylist.server/pkg/currency"
	"paylist.server/pkg/httpx"
	"paylist.server/pkg/httpx/httperr"
)

type convertResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}

type ratesResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (h *Handler) ConvertCurrencyHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	from := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("from")))
	to := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("to")))
	amountRaw := strings.TrimSpace(r.URL.Query().Get("amount"))

	if from == "" || to == "" || amountRaw == "" {
		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	amount, err := strconv.ParseFloat(amountRaw, 64)
	if err != nil || amount < 0 {
		return httperr.BadRequest(tr.TErr("error.invalid-amount"))
	}

	result, err := currency.Convert(ctx, from, to, amount)
	if err != nil {
		return httperr.BadRequest(err.Error())
	}

	httpx.HttpResponse(w, r, http.StatusOK, convertResponse{
		From:   from,
		To:     to,
		Amount: amount,
		Result: result,
	})

	return nil
}

func (h *Handler) GetCurrencyRatesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	base := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("base")))
	symbolsRaw := strings.TrimSpace(r.URL.Query().Get("symbols"))

	if base == "" {
		return httperr.BadRequest(tr.TErr("error.fields-not-filled"))
	}

	symbols := []string{}
	if symbolsRaw != "" {
		for _, symbol := range strings.Split(symbolsRaw, ",") {
			symbol = strings.ToUpper(strings.TrimSpace(symbol))
			if symbol != "" {
				symbols = append(symbols, symbol)
			}
		}
	}

	rates, err := currency.GetRates(ctx, base, symbols)
	if err != nil {
		return httperr.BadRequest(err.Error())
	}

	httpx.HttpResponse(w, r, http.StatusOK, ratesResponse{
		Base:  base,
		Rates: rates,
	})

	return nil
}

func (h *Handler) GetCurrenciesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	items, err := currency.GetCurrencies(ctx)
	if err != nil {
		return httperr.BadRequest(err.Error())
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, items)
	return nil
}
