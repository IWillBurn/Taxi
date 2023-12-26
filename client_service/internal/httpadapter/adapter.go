package httpadapter

import (
	"client_service/internal/config"
	"client_service/internal/httpadapter/requests"
	"client_service/internal/repo"
	"client_service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Adapter struct {
	Config             *config.HttpAdapterConfig
	DataBaseController repo.DataBaseController
	TripService        service.TripService
	server             *http.Server
}

func (a *Adapter) ListTrips(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	records, err := a.DataBaseController.GetTrips(userId)

	responseJson, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))
	return
}

func (a *Adapter) CreateTrip(w http.ResponseWriter, r *http.Request) {
	//userId := r.Header.Get("user_id")

	request := &requests.CreateTripRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tripId, err := a.TripService.CreateTrip(request.OfferId)

	response := make(map[string]string)
	response["trip_id"] = tripId
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (a *Adapter) GetTripById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	tripId := chi.URLParam(r, "trip_id")

	// userId := r.Header.Get("user_id")
	record, err := a.DataBaseController.GetTripByTripId(tripId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	responseJson, err := json.Marshal(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))
	return
}

func (a *Adapter) CancelTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	tripId := chi.URLParam(r, "trip_id")
	reason := chi.URLParam(r, "reason")

	//userId := r.Header.Get("user_id")
	err := a.TripService.CancelTrip(tripId, reason)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	return
}

func (a *Adapter) Serve() error {
	r := chi.NewRouter()

	apiRouter := chi.NewRouter()
	apiRouter.Get("/trips", a.ListTrips)
	apiRouter.Post("/trips", a.CreateTrip)
	apiRouter.Get("/trips/{trip_id}", a.GetTripById)
	apiRouter.Post("/trip/{trip_id}/cancel", a.CancelTrip)

	r.Mount(a.Config.BasePath, apiRouter)
	a.server = &http.Server{Addr: a.Config.ServeAddress, Handler: r}

	return a.server.ListenAndServe()
}

func (a *Adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(config *config.HttpAdapterConfig, DataBaseController repo.DataBaseController, tripService service.TripService) *Adapter {
	return &Adapter{
		Config:             config,
		DataBaseController: DataBaseController,
		TripService:        tripService,
	}
}
