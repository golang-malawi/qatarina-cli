package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/wakisa/qatarina-cli/internal/schema"
)

type InboxFilterParams struct{}

func FetchInbox(ctx context.Context, params *InboxFilterParams) ([]schema.TestCaseResponse, error) {
	path := "v1/me/test-cases/inbox"
	resp, err := Default().Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", string(bodyBytes))
	}

	var wrapper struct {
		TestCases []schema.TestCaseResponse `json:"test_cases"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.TestCases, nil
}
