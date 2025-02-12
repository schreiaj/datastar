package site

import (
	datastar "github.com/starfederation/datastar/sdk/go"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func setupHowTosUiIndicators(howTosUiIndicators chi.Router) error {

	howTosUiIndicators.Route("/ui-indicators", func(dataRouter chi.Router) {

		dataRouter.Post("/save", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		})

		dataRouter.Post("/optimistic", func(w http.ResponseWriter, r *http.Request) {
			signals := &OptimisticUpdateSignals{}
			datastar.ReadSignals(r, signals)
			sse := datastar.NewSSE(w, r)
			time.Sleep(2 * time.Second)
			sse.MergeSignals([]byte("{name:''}"))
			sse.MergeFragmentTempl(howTosUiIndicatorsUpdateSucess(signals.Name), datastar.WithSelector("#optimistic"), datastar.WithMergeAfter())
		})
	})

	return nil
}
