package picker

import (
	"context"

	"github.com/Kajmany/cata-up/cfg"
	"github.com/google/go-github/v63/github"
)

// TODO auth maybe
func GetClient() *github.Client {
	return github.NewClient(nil)
}

func GetRecentReleases(client *github.Client, source cfg.Source, page int, number int) ([]*github.RepositoryRelease, error) {
	ctx := context.Background() // TODO Am I making this right
	owner, repo := source.GetOwnerRepo()
	opt := &github.ListOptions{PerPage: number} // TODO use 'page' number
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, opt)
	// TODO might need to deal with response data
	return releases, err
}
