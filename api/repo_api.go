package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
	"github.com/google/go-github/v64/github"
)

// PlatformService provides helpers for calling platform_* endpoints.
type RepoService struct {
	client *client.FluencyClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewRepoService(client *client.FluencyClient) *RepoService {
	return &RepoService{client: client}
}

func (s *RepoService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *RepoService) ListRepos() (repos []*model.GithubRepo, err error) {
	req := &GenericDAORequest[model.GithubRepo]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.GithubRepo]
	if err := s.call("github_repo_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

type GitRepoContentRequest struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path"`
}

type GitRepoContentResponse struct {
	File      *github.RepositoryContent   `json:"file"`
	Directory []*github.RepositoryContent `json:"dirs"`
}

func (s *RepoService) GetRepoContent(repoId string, repoPath string) (resp *GitRepoContentResponse, err error) {
	req := &GitRepoContentRequest{
		ID:   repoId,
		Path: repoPath,
	}
	if err := s.call("github_repo_get_content", req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type ImportRepoObjectsRequest struct {
	ID    string   `json:"id,omitempty"`
	Paths []string `json:"paths"`
}

func (s *RepoService) ImportRepoProcessors(repoId string, filePaths []string) (err error) {
	req := &ImportRepoObjectsRequest{
		ID:    repoId,
		Paths: filePaths,
	}
	if err := s.call("platform_import_processors", req, nil); err != nil {
		return err
	}
	return nil
}

func (s *RepoService) ImportAppTemplates(repoId string, filePaths []string) (err error) {
	req := &ImportRepoObjectsRequest{
		ID:    repoId,
		Paths: filePaths,
	}
	if err := s.call("platform_import_application_templates", req, nil); err != nil {
		return err
	}
	return nil
}
