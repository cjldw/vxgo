package vxgo

import (
	"bufio"
	"bytes"
	"errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsClonedRepo() bool {
	repoPath := filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName)
	_, err := os.Stat(repoPath)
	if err != nil {
		return false
	}
	return true
}

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
		log.Printf("can't open code repository failure: %v\n", err)
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

/// go-git package use don't happy

func GitCloneCmd() (bool, error) {
	gitRepoDir := filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName)
	_, err := os.Stat(gitRepoDir)
	if err == nil {
		log.Printf("gitRepoDir is clone : %v\n", err)
		return false, errors.New("gitRepoDir: " + gitRepoDir + " exists")
	}
	cmd := exec.Command("git", "clone", VxCfg.GitRepo)
	cmd.Dir = VxCfg.WorkDir
	resultBytes, err := cmd.Output()
	if err != nil {
		return false, err
	}
	log.Printf("git clone " + VxCfg.GitRepo + " status: " + string(resultBytes) + "\n")
	return true, nil
}

func GitPullCmd() (bool, error) {
	gitRepoDir := filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName)
	cmd := exec.Command("git", "pull")
	cmd.Dir = gitRepoDir
	resultBytes, err := cmd.Output()
	log.Printf("run git pull output: %s\n", string(resultBytes))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GitShowCaseCmd() (*NewCommitPoint, error) {
	gitRepoDir := filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName)
	cmd := exec.Command("git", "show", "-n", "1", "--name-only", "--pretty=%H")
	cmd.Dir = gitRepoDir
	rsBytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lineNo := 0
	scan := bufio.NewScanner(bytes.NewBuffer(rsBytes))
	commit := new(NewCommitPoint)
	for scan.Scan() {
		lineNo++
		txt := strings.Trim(strings.TrimSpace(scan.Text()), "\"")
		if len(txt) <= 0 {
			continue
		}
		if lineNo <= 1 {
			commit.CommitID = txt
			continue
		}
		commit.Files = append(commit.Files, filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName, txt))
	}
	return commit, nil
}
