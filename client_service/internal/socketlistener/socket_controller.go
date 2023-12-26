package socketlistener

import (
	"client_service/internal/config"
	"client_service/internal/models"
	"client_service/internal/service"
	"client_service/internal/socketlistener/listeners"
	"client_service/internal/socketlistener/publishers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

type SocketController struct {
	Config     *config.SocketConfig
	Clients    *sync.Map
	Connector  *Connector
	Publishers *map[string]*publishers.Publisher
	server     *http.Server
}

func (s *SocketController) RegisterNewPublisher(key string) *listeners.StatusListener {

	publisher := &publishers.Publisher{
		Subscribers: sync.Map{},
		Broadcast:   make(chan models.SocketMessage),
		Key:         key,
	}

	listener := &listeners.StatusListener{
		Clients:   s.Clients,
		Publisher: publisher,
	}

	return listener
}

func (s *SocketController) Serve() error {
	r := chi.NewRouter()

	apiRouter := chi.NewRouter()
	apiRouter.Get("/", s.Connector.HandleConnections)

	r.Mount(s.Config.BasePath, apiRouter)
	fmt.Println(s.Config.ServeAddress)
	fmt.Println(s.Config.BasePath)
	s.server = &http.Server{Addr: s.Config.ServeAddress, Handler: r}

	for _, pub := range *s.Publishers {
		listener := listeners.StatusListener{
			Clients:   s.Clients,
			Publisher: pub,
		}
		go listener.HandleMessages()
	}

	return s.server.ListenAndServe()
}

func NewSocketController(config *config.SocketConfig, tripService service.TripService) (*SocketController, error) {
	clients := sync.Map{}
	publisher := publishers.NewPublisher("status")
	publishers := make(map[string]*(publishers.Publisher))
	publishers["status"] = publisher
	return &SocketController{
		Config:     config,
		Clients:    &clients,
		Publishers: &publishers,
		Connector: &Connector{
			TripService: tripService,
			Clients:     &clients,
			Publishers:  &publishers,
		},
	}, nil
}
