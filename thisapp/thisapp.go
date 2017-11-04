package thisapp

import (
	"log"
	"expvar"
	"net/http"
	"time"
	"os"
)

type AppModule struct {
	config     *Configs
	stats      *expvar.Int
	log        *log.Logger
	httpClient http.Client
	ProductDB  ProductDB
}

func NewAppModule(env string) (*AppModule) {
	config, err := readConfig(env)
	if err != nil {
		log.Fatalf("Failed to read config because %s", err.Error())
	}

	return newAppModuleWithConfig(&config)
}

func newAppModuleWithConfig(config *Configs) (*AppModule) {
	productDB, err := NewProductDB(config)
	if err != nil {
		log.Fatalf("Failed to create product DB because %s", err.Error())
	}

	return &AppModule{
		config:     config,
		stats:      expvar.NewInt("requestAPI"),
		log:        log.New(os.Stdout, "debug:", log.Ldate|log.Ltime|log.Lshortfile),
		httpClient: *newHttpClient(config),
		ProductDB:  productDB,
	}
}

func newHttpClient(cfg *Configs) *http.Client {
	client := http.Client{}
	client.Timeout = time.Duration(10 * time.Second)
	return &client
}
