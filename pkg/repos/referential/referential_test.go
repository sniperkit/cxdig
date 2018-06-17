package referential

import (
	"testing"

	"github.com/sniperkit/cxdig/pkg/config"
	"github.com/sniperkit/cxdig/pkg/types"

	"github.com/stretchr/testify/assert"
)

/*
func TestSaveReferentialToJSON(t *testing.T) {
	id := core.ProjectID("testJSON")
	diff := types.ProjectReferential{
		types.LocalFile{ID: 1},
		types.LocalFile{ID: 2},
		types.LocalFile{ID: 3},
	}
	assert.NoError(t, SaveReferentialToJSON(id, diff))
	_, err := os.Stat("./" + id.String() + ".[referential].json")
	assert.NoError(t, err)
	os.Remove("./" + id.String() + ".[referential].json")
}*/

func TestBuildProjectReferential(t *testing.T) {
	commits := []types.CommitInfo{
		{
			Number:   1,
			CommitID: "shaCommit1",
			Author: types.AuthorInfo{
				Name:  "author1",
				Email: "author1@mail.com",
			},
			Message: "message 1",
			Changes: []types.FileChange{
				{
					Type:     types.FileChangeModified,
					FilePath: "first/path",
				},
				{
					Type:     types.FileChangeDeleted,
					FilePath: "second/path",
				},
			},
		},
		{
			Number:   2,
			CommitID: "shaCommit2",
			Author: types.AuthorInfo{
				Name:  "author2",
				Email: "author2@mail.com",
			},
			Message: "message 2",
			Changes: []types.FileChange{
				{
					Type:     types.FileChangeModified,
					FilePath: "first/path",
				},
				{
					Type:     types.FileChangeDeleted,
					FilePath: "second/path",
				},
			},
		},
		{
			Number:   3,
			CommitID: "shaCommit3",
			Author: types.AuthorInfo{
				Name:  "author3",
				Email: "author3@mail.com",
			},
			Message: "message 3",
			Changes: []types.FileChange{
				{
					Type:     types.FileChangeModified,
					FilePath: "first/path",
				},
				{
					Type:     types.FileChangeDeleted,
					FilePath: "",
				},
			},
		},
	}
	assert.NotPanics(t, func() { BuildProjectReferential(commits, config.NewFileTypeRegistry()) })
	commits[2] = types.CommitInfo{
		Number:   3,
		CommitID: "shaCommit3",
		Author: types.AuthorInfo{
			Name:  "author3",
			Email: "author3@mail.com",
		},
		Message: "message 3",
		Changes: []types.FileChange{
			{
				Type:     types.FileChangeModified,
				FilePath: "first/path",
			},
			{
				Type:     types.FileChangeType("Unknown"),
				FilePath: "second/path",
			},
		},
	}
	assert.Panics(t, func() { BuildProjectReferential(commits, config.NewFileTypeRegistry()) })
}

func TestNext(t *testing.T) {
	id := LocalDirectoryID(1)
	id = id.Next()
	assert.Equal(t, LocalDirectoryID(2), id)
}
