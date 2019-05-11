package diplomat

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type OutputDirectoryTestSuite struct {
	suite.Suite
}

func (o OutputDirectoryTestSuite) Cleanup(testDir string) {
	o.Require().NoError(os.RemoveAll(testDir))
}

func (o OutputDirectoryTestSuite) TestWriteFile() {
	testRoot := filepath.Join(os.TempDir(),uuid.New().String())
	err := os.MkdirAll(testRoot, DefaultDirectoryPerm)
	o.Require().NoError(err)
	defer o.Cleanup(testRoot)
	od := NewOutputDirectory(testRoot)
	data := []byte("content of helloworld")
	o.Require().NoError(od.WriteFile("helloworld", data))
	expectedFilePath := filepath.Join(testRoot,"helloworld")
	o.Require().FileExists(expectedFilePath)
	content, err := ioutil.ReadFile(expectedFilePath)
	o.Require().NoError(err)
	o.Equal(content, data)
}

func TestOutputDirectory(t *testing.T) {
	suite.Run(t, &OutputDirectoryTestSuite{})
}

