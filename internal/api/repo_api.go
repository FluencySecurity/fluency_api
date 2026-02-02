package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
)

func getParserType(processor_type string) (string, error) {
	switch processor_type {
	case "fpl_processor":
		return "parsers", nil
	case "fpl_receiver":
		return "receivers", nil
	case "fpl_packer":
		return "packers", nil
	default:
		return "", fmt.Errorf("unknown processor type: %s", processor_type)
	}
}

func (c *Client) ImportProcessor(processor_type string, repoName string) (err error) {

	repoService := fluencyAPI.NewRepoService(c.fluencyClient)

	resp, err := repoService.ListRepos()

	if err != nil {
		c.Logger.Error("failed to list repos", "error", err)
		return fmt.Errorf("failed to list repos: %w", err)
	}
	var repoID string

	if repoName == "" {
		if len(resp) == 0 {
			c.Logger.Error("no repos found")
			return fmt.Errorf("no repos found")
		}
		repoID = resp[0].ID
	} else {
		for _, repo := range resp {
			if repo.Repo == repoName {
				repoID = repo.ID
				break
			}
		}
	}
	if repoID == "" {
		c.Logger.Error("repo not found", "repo", repoName)
		return fmt.Errorf("repo not found: %s", repoName)
	}

	parserType, err := getParserType(processor_type)
	if err != nil {
		c.Logger.Error("failed to get parser type", "error", err)
		return fmt.Errorf("unknown parser type %s: %w", processor_type, err)
	}

	metaPath := fmt.Sprintf("%s/meta", parserType)

	contentResp, err := repoService.GetRepoContent(repoID, metaPath)
	if err != nil {
		c.Logger.Error("failed to get repo content", "error", err)
		return fmt.Errorf("failed to get repo content: %w", err)
	}

	var files []string
	for _, file := range contentResp.Directory {
		if file.GetType() == "file" {
			files = append(files, *file.Path)
		}
	}
	err = repoService.ImportRepoProcessors(repoID, files)
	if err != nil {
		c.Logger.Error("failed to import repo processors", "error", err)
		return fmt.Errorf("failed to import repo processors: %w", err)
	}
	return nil
}

func (c *Client) ImportAppTemplate(repoName string) (err error) {

	repoService := fluencyAPI.NewRepoService(c.fluencyClient)

	resp, err := repoService.ListRepos()

	if err != nil {
		c.Logger.Error("failed to list repos", "error", err)
		return fmt.Errorf("failed to list repos: %w", err)
	}
	var repoID string

	if repoName == "" {
		if len(resp) == 0 {
			c.Logger.Error("no repos found")
			return fmt.Errorf("no repos found")
		}
		repoID = resp[0].ID
	} else {
		for _, repo := range resp {
			if repo.Repo == repoName {
				repoID = repo.ID
				break
			}
		}
	}
	if repoID == "" {
		c.Logger.Error("repo not found", "repo", repoName)
		return fmt.Errorf("repo not found: %s", repoName)
	}

	//metaPath := fmt.Sprintf("applications/meta", parserType)

	contentResp, err := repoService.GetRepoContent(repoID, "applications/meta")
	if err != nil {
		c.Logger.Error("failed to get repo content", "error", err)
		return fmt.Errorf("failed to get repo content: %w", err)
	}

	var files []string
	for _, file := range contentResp.Directory {
		if file.GetType() == "file" {
			files = append(files, *file.Path)
		}
	}
	err = repoService.ImportAppTemplates(repoID, files)
	if err != nil {
		c.Logger.Error("failed to import repo processors", "error", err)
		return fmt.Errorf("failed to import repo processors: %w", err)
	}
	return nil
}
