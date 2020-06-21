package remote

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/google/go-github/v32/github"
	"github.com/hatzelencio/branch-protection/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	token       = "GITHUB_TOKEN"
	ghrepo      = "GITHUB_REPOSITORY"
	ghworkspace = "GITHUB_WORKSPACE"
	path        = "INPUT_PATH"
)

var (
	ctx    context.Context
	cli    GithubClient
	writer io.Writer
)

func init() {
	writer = os.Stderr

	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(token)},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli = NewGithubClient(tc, nil)
}

// GithubGitService interface
type GithubGitService interface {
	UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *github.ProtectionRequest) (*github.Protection, *github.Response, error)
}

// GithubClient is a wrapper of github.Client
type GithubClient struct {
	Repositories GithubGitService
	*github.Client
}

// NewGithubClient Create a new github client
func NewGithubClient(client *http.Client, repoMock GithubGitService) GithubClient {
	if repoMock != nil {
		return GithubClient{
			Repositories: repoMock,
		}
	}

	cli := github.NewClient(client)
	return GithubClient{
		Repositories: cli.Repositories,
	}
}

// ValidateInputs validate if GITHUB_TOKEN and INPUT_REFS are present like environment variables
func ValidateInputs() error {
	if len(os.Getenv(token)) == 0 {
		return fmt.Errorf("%v is required env variable to trigger this action", token)
	}

	path := getBranchProtectionConfigPath()
	exists, _ := utils.FileExists(path)
	if !exists {
		return fmt.Errorf("ensure if config file %v exists", path)
	}

	return nil
}

func getOwnerRepo() (string, string) {
	ownerRepo := strings.Split(os.Getenv(ghrepo), "/")
	return ownerRepo[0], ownerRepo[1]
}

func getBranchProtectionConfigPath() string {
	var prefix = ""

	path := os.Getenv(path)

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "C:\\") {
		prefix = fmt.Sprintf("%v/", os.Getenv(ghworkspace))
	}

	return fmt.Sprintf("%v%v", prefix, path)
}

// BranchProtection struct
type BranchProtection struct {
	Branch     string                    `json:"branch"`
	Protection *github.ProtectionRequest `json:"protection"`
}

func getBranchProtectionRequests() ([]BranchProtection, error) {
	var protections []BranchProtection
	data, err := ioutil.ReadFile(getBranchProtectionConfigPath())

	if err != nil {
		return nil, err
	}

	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &protections)

	if err != nil {
		return nil, err
	}

	return protections, nil
}

// UpdateBranchProtection update a protection by branch
func UpdateBranchProtection() error {
	var wg sync.WaitGroup
	requests, err := getBranchProtectionRequests()

	if err != nil {
		return err
	}

	wg.Add(len(requests))
	owner, repo := getOwnerRepo()

	for _, bp := range requests {
		go func(bp BranchProtection) {
			defer wg.Done()
			_, _, err := cli.Repositories.UpdateBranchProtection(ctx, owner, repo, bp.Branch, bp.Protection)

			if err != nil {
				log.Fatal(err)
			}

			_, err = fmt.Fprintln(writer, fmt.Sprintf("branch %v has been protected", bp.Branch))

			if err != nil {
				log.Fatal(err)
			}
		}(bp)
	}

	wg.Wait()

	return nil
}
