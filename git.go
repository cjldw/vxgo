package vxgo

import (
	"errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"path/filepath"
	"strings"
)

func CloneRepo() error {
	r, err := git.PlainClone(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName), false, &git.CloneOptions{
		URL: VxCfg.GitRepo,
	})
	if err != nil {
		log.Printf("git clone :%s failure: %v\n", VxCfg.GitRepo, err)
		return err
	}
	ref, err := r.Head()
	if err != nil {
		log.Printf("git check HEAD failure: %v\n", err)
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Printf("read commit failure: %v\n", err)
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		log.Printf("commit tree failure: %v\n", err)
		return err
	}
	err = tree.Files().ForEach(func(file *object.File) error {
		log.Println(file.Name)
		return nil
	})
	if err != nil {
		log.Printf("read tree file failure: %v\n", err)
		return err
	}
	return nil
}

func PullRepo() error {
	r, err := git.PlainOpen(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName))
	if err != nil {
		log.Printf("open git repository failure: %v\n", err)
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		log.Fatalf("repository worktree failure: %v\n", err)
		return err
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		log.Printf("git pull repository failure: %v\n", err)
		return err
	}
	return nil
}

// GitShowCase get latest directory [source/_posts] files same with
// `git log --number 1 --pretty=format:%h --name-only`
func GitShowCase() (*NewCommitPoint, error) {
	r, err := git.PlainOpen(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName))
	if err != nil {
		log.Printf("open git repository failure: %v\n", err)
		return nil, err
	}
	ref, err := r.Head()
	if err != nil {
		log.Printf("git checkout HEAD failure: %v\n", err)
		return nil, err
	}
	logsInfo, err := r.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	if err != nil {
		log.Printf("git commit log failure: %v\n", err)
		return nil, err
	}
	commit, err := logsInfo.Next()
	if err != nil {
		log.Printf("git fetch next commit point failure: %v\n", err)
		return nil, err
	}
	fileState, err := commit.Stats()
	if err != nil {
		return nil, err
	}
	files := strings.Split(fileState.String(), "\n")
	fileLength := len(files)
	if fileLength <= 0 {
		return nil, errors.New("this commit no blog file change")
	}
	commitPoint := &NewCommitPoint{
		CommitID: commit.ID().String(),
		Files:    []string{},
	}
	for i := 0; i < fileLength; i++ {
		fileName := strings.Trim(strings.Split(files[i], "|")[0], " ")
		if !strings.HasPrefix(fileName, "source/_posts") {
			continue
		}
		commitPoint.Files = append(commitPoint.Files, filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName, fileName))
	}

	return commitPoint, nil
}
