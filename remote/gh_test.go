package remote

import (
	"bytes"
	"fmt"
	"github.com/google/go-github/v39/github"
	"github.com/hatzelencio/branch-protection/utils/mocks"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

type envVariables struct {
	token      string
	repository string
	path       string
}

func init() {
	cli = NewGithubClient(nil, &mocks.MockClient{})
}

func TestUpdateBranchProtectionSuccess(t *testing.T) {
	var buf bytes.Buffer
	configFile := createTmpBranchProtectionConfigFile()

	defer func() {
		os.Remove(configFile.Name())
		writer = os.Stderr
	}()

	writer = &buf

	_, err := configFile.Write([]byte(`
- branch: test-branch
  protection:
    required_pull_request_reviews:
      required_approving_review_count: 1
    allow_deletions: true`))

	if err != nil {
		t.Fatal(fmt.Sprintf("not possible read a config file %v", configFile.Name()))
	}

	params := envVariables{
		token:      "secret-token",
		repository: "owner/owner-repo",
		path:       configFile.Name(),
	}

	setEnvVariables(params)
	mockUpdateBranchSuccess()

	err = UpdateBranchProtection()

	if err != nil {
		t.Fatal("it can'n be create a protection")
	}

	expected := `branch test-branch has been protected
`
	if buf.String() != expected {
		t.Fatal("contest are not equals")
	}
}

func setEnvVariables(env envVariables) {
	_ = os.Setenv("GITHUB_TOKEN", env.token)
	_ = os.Setenv("GITHUB_REPOSITORY", env.repository)
	_ = os.Setenv("INPUT_PATH", env.path)
}

func createTmpBranchProtectionConfigFile() (file *os.File) {
	file, err := ioutil.TempFile("", "branch_protection_*.yml")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func mockUpdateBranchSuccess() {
	mocks.UpdateBranchProtectionFunc = func(ctx context.Context, owner, repo, branch string, preq *github.ProtectionRequest) (*github.Protection, *github.Response, error) {
		res := github.Response{
			Response: &http.Response{StatusCode: 200},
		}
		return &github.Protection{}, &res, nil
	}
}
