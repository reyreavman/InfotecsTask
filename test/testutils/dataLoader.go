package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type DataLoader struct {
	basePath string
}

func NewDataLoader() *DataLoader {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return &DataLoader{
		basePath: filepath.Join(dir, "..", "testdata"),
	}
}

func (dl *DataLoader) LoadJSONFixture(fixturePath string, out any) error {
	fullPath := filepath.Join(dl.basePath, "fixtures", "json", fixturePath)
	data, err := ioutil.ReadFile(fullPath)

	if err != nil {
		return fmt.Errorf("failed to read JSON fixture: %w", err)
	}

	return json.Unmarshal(data, out)
}
