package vxgo

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"path/filepath"
)

func CloneRepo() {
	r, err := git.PlainClone(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName), false, &git.CloneOptions{
		URL: VxCfg.GitRepo,
	})
	if err != nil {
		log.Fatalf("git clone :%s failure: %v\n", VxCfg.GitRepo, err)
	}
	ref, err := r.Head()
	if err != nil {
		log.Fatalf("git check HEAD failure: %v\n", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("read commit failure: %v\n", err)
	}
	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("commit tree failure: %v\n", err)
	}
	err = tree.Files().ForEach(func(file *object.File) error {
		log.Println(file.Name)
		return nil
	})
	if err != nil {
		log.Fatalf("read tree file failure: %v\n", err)
	}
}

func PullRepo() {

	r, err := git.PlainOpen(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName))
	if err != nil {
		log.Fatalf("open git repository failure: %v\n", err)
	}
	wt, err := r.Worktree()
	if err != nil {
		log.Fatalf("repository worktree failure: %v\n", err)
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		log.Fatalf("git pull repository failure: %v\n", err)
	}
	ref, err := r.Head()
	if err != nil {
		log.Fatalf("git checkout HEAD failure: %v\n", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("git commit log failure: %v\n", err)
	}
	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("git get tree failure: %v\n", err)
	}
	err = tree.Files().ForEach(func(file *object.File) error {
		log.Println("===================")
		log.Println(file.Name)
		log.Println("===================")
		return nil
	})
	if err != nil {
		log.Fatalf("git worktree files failure: %v\n", err)
	}
}

func GitShowCase() {
	r, err := git.PlainOpen(filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName))
	if err != nil {
		log.Fatalf("open git repository failure: %v\n", err)
	}
	ref, err := r.Head()
	if err != nil {
		log.Fatalf("git checkout HEAD failure: %v\n", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("git commit log failure: %v\n", err)
	}
	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("git get tree failure: %v\n", err)
	}
	err = tree.Files().ForEach(func(file *object.File) error {
		log.Println("===================")
		log.Println(file.Name)
		log.Println("===================")
		return nil
	})
	if err != nil {
		log.Fatalf("git worktree files failure: %v\n", err)
	}
}
