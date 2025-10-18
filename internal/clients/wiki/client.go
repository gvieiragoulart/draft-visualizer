package esports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gvieiragoulart/draft-visualizer/internal/clients"
	"github.com/gvieiragoulart/draft-visualizer/internal/clients/esports/dto"
	mwclient "cgt.name/pkg/go-mwclient"
)

type WikiClient struct {
	mwclient.Client
}

func NewClient(apiKey string) *WikiClient {
	w, err := mwclient.New("https://lol.fandom.com/api.php")
	if err != nil {
		log.Fatalf("Error creating wiki client: %v", err)
	}

	return &WikiClient{
		mwclient.Client: w,
	}
}

func (c *WikiClient) Login(username, password string) error {
	return c.Client.Login(username, password)
}

func (c *WikiClient) GetPage(page string) (string, error) {
	return c.Client.GetPage(page)
}