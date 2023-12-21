package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"offering_service/internal/config"
	"offering_service/internal/httpadapter/requests"
	"offering_service/internal/httpadapter/responses"
	"offering_service/internal/service"
)

type Adapter struct {
	Config          *config.HttpAdapterConfig
	OfferingService service.OfferingService
	SigningService  service.SigningService
	server          *http.Server
}

func (a *Adapter) CreateOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	request := &requests.CreateOfferRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dataObj := a.OfferingService.GetPrice(*request)
	id, err := a.SigningService.Encode(dataObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.OfferResponse{Id: id, From: dataObj.From, To: dataObj.To, ClientId: dataObj.ClientId, Price: dataObj.Price}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))
	return
}

func (a *Adapter) ParseOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	offerId := chi.URLParam(r, "offer_id")

	data, err := a.SigningService.Decode(offerId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	response := responses.OfferResponse{Id: offerId, From: data.From, To: data.To, Price: data.Price, ClientId: data.ClientId}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(responseJson))
	return
}

func (a *Adapter) Serve() error {
	r := chi.NewRouter()

	apiRouter := chi.NewRouter()
	apiRouter.Post("/offers", a.CreateOffer)
	apiRouter.Get("/offers/{offer_id}", a.ParseOffer)

	r.Mount(a.Config.BasePath, apiRouter)
	fmt.Println(a.Config.BasePath)
	fmt.Println(a.Config.ServeAddress)
	a.server = &http.Server{Addr: a.Config.ServeAddress, Handler: r}

	return a.server.ListenAndServe()
}

func (a *Adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(config *config.HttpAdapterConfig, offeringService service.OfferingService, signingService service.SigningService) *Adapter {
	return &Adapter{
		Config:          config,
		OfferingService: offeringService,
		SigningService:  signingService,
	}
}
