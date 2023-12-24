package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"offering_service/internal/config"
	"offering_service/internal/httpadapter/requests"
	"offering_service/internal/service"
)

type Adapter struct {
	Config          *config.HttpAdapterConfig
	DatabaseService service.DataBaseService
	TripService     service.TripService
	server          *http.Server
}

func (a *Adapter) ListTrips(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")
	records := a.DatabaseService.GetTrips(userId)

	responseJson, err := json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))
	return
}

func (a *Adapter) CreateTrip(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("user_id")

	request := &requests.CreateTripRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.TripService.CreateTrip(userId, request.OfferId)

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
	record, err := a.DatabaseService.GetTripById(tripId)

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
	// reason := chi.URLParam(r, "reason")

	userId := r.Header.Get("user_id")
	err := a.TripService.CancelTrip(userId, tripId)

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
	fmt.Println(a.Config.BasePath)
	fmt.Println(a.Config.ServeAddress)
	a.server = &http.Server{Addr: a.Config.ServeAddress, Handler: r}

	return a.server.ListenAndServe()
}

func (a *Adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(config *config.HttpAdapterConfig, databaseService service.DataBaseService, tripService service.TripService) *Adapter {
	return &Adapter{
		Config:          config,
		DatabaseService: databaseService,
		TripService:     tripService,
	}
}
