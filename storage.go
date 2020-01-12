package vxgo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Storager interface {
	SaveCommit(commit string) (bool, error)
	QueryCommit(commit string) bool
}

type Dumper struct {
	dumpfile string
	isLoad   bool
	dumpData map[string]interface{}
	dumpLock *sync.RWMutex
}

var (
	drivers = &sync.Map{}
)

func GetDumper(filename ...string) *Dumper {
	_, file, _, _ := runtime.Caller(0)
	dumpfile := filepath.Join(filepath.Dir(file), "dumpfile.json")
	if len(filename) > 0 {
		dumpfile = filename[0]
	}
	item, ok := drivers.Load(dumpfile)
	if ok {
		return item.(*Dumper)
	}
	dump := &Dumper{
		dumpfile: dumpfile,
		dumpData: make(map[string]interface{}),
		dumpLock: &sync.RWMutex{},
	}
	drivers.Store(dumpfile, dump)
	return dump
}

func (dp *Dumper) load() (bool, error) {
	dp.dumpLock.Lock()
	defer dp.dumpLock.Unlock()
	if dp.isLoad {
		return true, nil
	}
	file, err := os.OpenFile(dp.dumpfile, os.O_RDWR|os.O_CREATE, 644)
	if err != nil {
		log.Printf("open file :%s failure %v\n", dp.dumpfile, err)
		return false, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Printf("read file :%s failure %v\n", dp.dumpfile, err)
		return false, err
	}
	err = json.Unmarshal(fileBytes, &dp.dumpData)
	if err != nil {
		log.Printf("json unmarshal dump file failure (maybe file not found): %v\n", err)
	}
	dp.isLoad = true
	return true, nil
}

func (dp *Dumper) save() (bool, error) {
	dp.dumpLock.Lock()
	defer dp.dumpLock.Unlock()
	fileBytes, err := json.MarshalIndent(dp.dumpData, "", "    ")
	if err != nil {
		log.Printf("json marshal failure %v\n", err)
		return false, err
	}
	err = ioutil.WriteFile(dp.dumpfile, fileBytes, 644)
	if err != nil {
		log.Printf("write to file failure %v\n", err)
		return false, err
	}
	return true, nil
}

func (dp *Dumper) SaveCommit(commit string, value interface{}) (bool, error) {
	_, err := dp.load()
	if err != nil {
		return false, err
	}
	dp.dumpData[commit] = value
	return dp.save()
}

func (dp *Dumper) QueryCommit(commit string) (interface{}, bool, error) {
	_, err := dp.load()
	if err != nil {
		return nil, false, err
	}
	value, ok := dp.dumpData[commit]
	if !ok {
		return nil, false, errors.New("item not storage")
	}
	return value, true, nil
}
