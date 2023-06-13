package clientset

import (
	"context"
	"fmt"
	"strings"

	"github.com/coralogix/coralogix-operator/controllers/clientset/rest"
)

type TCOPolicies struct {
	client *rest.Client
}

func (t TCOPolicies) CreateTCOPolicy(ctx context.Context, jsonContent string) (string, error) {
	return t.client.Post(ctx, "/policies", "application/json", jsonContent)
}

func (t TCOPolicies) GetTCOPolicy(ctx context.Context, id string) (string, error) {
	return t.client.Get(ctx, fmt.Sprintf("/policies/%s", id))
}

func (t TCOPolicies) UpdateTCOPolicy(ctx context.Context, id string, jsonContent string) (string, error) {
	return t.client.Put(ctx, fmt.Sprintf("/policies/%s", id), "application/json", jsonContent)
}

func (t TCOPolicies) DeleteTCOPolicy(ctx context.Context, id string) error {
	_, err := t.client.Delete(ctx, fmt.Sprintf("/policies/%s", id))
	return err
}

func NewTCOPoliciesClient(c *CallPropertiesCreator) *TCOPolicies {
	targetUrl := "https://" + strings.Replace(c.targetUrl, "ng-api-grpc", "api", 1) + "/api/v1/external/tco"
	client := rest.NewRestClient(targetUrl, c.apiKey)
	return &TCOPolicies{client: client}
}
