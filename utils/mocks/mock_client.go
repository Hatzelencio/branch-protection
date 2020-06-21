package mocks

import (
	"github.com/google/go-github/v32/github"
	"golang.org/x/net/context"
)

// MockClient mock of github.Client
type MockClient struct {}

var (
	// UpdateBranchProtectionFunc mock variable to create a git branch protection to github repository
	UpdateBranchProtectionFunc func(ctx context.Context, owner, repo, branch string, preq *github.ProtectionRequest) (*github.Protection, *github.Response, error)
)

// UpdateBranchProtection mock variable to create a git branch protection to github repository
func (m *MockClient) UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *github.ProtectionRequest) (*github.Protection, *github.Response, error){
	return UpdateBranchProtectionFunc(ctx, owner, repo, branch, preq)
}