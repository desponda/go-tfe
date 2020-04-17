package tfe

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AdminRuns = (*adminRuns)(nil)

// AdminRuns Users Admin API contains endpoints to help site administrators manage
// user accounts.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/admin runs.html
type AdminRuns interface {
	// List all the runs of the given installation.
	List(ctx context.Context, options AdminRunsListOptions) (*AdminRunsList, error)
}

// runs implements Users.
type adminRuns struct {
	client *Client
}

// AdminRunsList represents a list of runs.
type AdminRunsList struct {
	*Pagination
	Items []*Run
}

// AdminRunsListOptions represents the options for listing runs.
type AdminRunsListOptions struct {
	ListOptions
}

// List all the runs of the terraform enterprise installation.
func (s *adminRuns) List(ctx context.Context, options AdminRunsListOptions) (*AdminRunsList, error) {
	u := fmt.Sprintf("admin/runs")
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rl := &AdminRunsList{}
	err = s.client.do(ctx, req, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// ForceCancel is used to forcefully cancel a run by its ID.
func (s *adminRuns) ForceCancel(ctx context.Context, runID string, options RunForceCancelOptions) error {
	if !validStringID(&runID) {
		return errors.New("invalid value for run ID")
	}

	u := fmt.Sprintf("admin/runs/%s/actions/force-cancel", url.QueryEscape(runID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
