package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const graphqlEndpoint = "https://api.github.com/graphql"

// GraphQLError represents an error from the GraphQL API.
type GraphQLError struct {
	Message string `json:"message"`
}

// QueryGraphQL executes a GraphQL query against GitHub's API and unmarshals
// the response into result. The result should be a pointer to a struct with
// a Data field matching the expected query shape.
func QueryGraphQL(ctx context.Context, token, query string, variables map[string]interface{}, result interface{}) error {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal graphql: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, graphqlEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create graphql request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("graphql request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read graphql response: %w", err)
	}

	// Check for GraphQL-level errors first
	var errCheck struct {
		Errors []GraphQLError `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &errCheck); err == nil && len(errCheck.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", errCheck.Errors[0].Message)
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("decode graphql response: %w", err)
	}
	return nil
}
