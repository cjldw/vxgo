package vxgo

import (
	"log"
	"os/exec"
	"path/filepath"
)

func HexoDeploy() (string, error) {
	cmd := exec.Command(VxCfg.HexoBin, "generate", "-d")
	cmd.Dir = filepath.Join(VxCfg.WorkDir, VxCfg.GitRepoName)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("run hexo generate -d failure: %v\n", err)
		return "", err
	}
	log.Printf("hexo generate -d output: %s\n", string(output))
	return string(output), nil
}
