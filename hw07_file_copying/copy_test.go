package main

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CopyDataSuite struct {
	suite.Suite
	f *os.File
	b []byte
}

func TestCopyDataSuite(t *testing.T) {
	suite.Run(t, new(CopyDataSuite))
}

func generateFilename() string {
	return strings.Join([]string{os.TempDir(), "otus_out" + strconv.Itoa(rand.Int())}, string(os.PathSeparator))
}

func (c *CopyDataSuite) SetupSuite() {
	var err error
	c.f, err = ioutil.TempFile("", "otus_copy")
	c.Require().NoError(err)

	bigBuff := make([]byte, 2<<9)
	_, err = rand.Read(bigBuff)
	c.Require().NoError(err)

	_, err = c.f.Write(bigBuff)
	c.Require().NoError(err)

	c.b = bigBuff

	err = c.f.Close()
	c.Require().NoError(err)

	rand.Seed(time.Now().Unix())
}

func (c *CopyDataSuite) TearDownSuite() {
	err := os.Remove(c.f.Name())
	c.Require().NoError(err)
}

func (c *CopyDataSuite) TestOffsetExceeded() {
	err := Copy(c.f.Name(), "", 2<<9+1, 0)
	c.Require().ErrorIs(err, ErrOffsetExceedsFileSize)
}

func (c *CopyDataSuite) TestCopyPartFile() {
	tmpName := generateFilename()
	err := Copy(c.f.Name(), tmpName, 0, 10)
	c.Require().NoError(err)

	f2, err := os.Open(tmpName)
	c.Require().NoError(err)
	data2, err := ioutil.ReadAll(f2)
	c.Require().NoError(err)

	c.Require().True(bytes.Equal(c.b[:10], data2))

	err = os.Remove(tmpName)
	c.Require().NoError(err)
}

func (c *CopyDataSuite) TestCopyFileOffset() {
	tmpName := generateFilename()
	err := Copy(c.f.Name(), tmpName, 10, 10)
	c.Require().NoError(err)

	f2, err := os.Open(tmpName)
	c.Require().NoError(err)
	data2, err := ioutil.ReadAll(f2)
	c.Require().NoError(err)

	c.Require().True(bytes.Equal(c.b[10:20], data2))

	err = os.Remove(tmpName)
	c.Require().NoError(err)
}

func (c *CopyDataSuite) TestCopyFullFile() {
	tmpName := generateFilename()
	err := Copy(c.f.Name(), tmpName, 0, 0)
	c.Require().NoError(err)

	f2, err := os.Open(tmpName)
	c.Require().NoError(err)
	data2, err := ioutil.ReadAll(f2)
	c.Require().NoError(err)

	c.Require().True(bytes.Equal(c.b, data2))

	err = os.Remove(tmpName)
	c.Require().NoError(err)
}

func (c *CopyDataSuite) TestErrorOffset() {
	tmpName := generateFilename()
	err := Copy(c.f.Name(), tmpName, -1, 0)
	c.Require().ErrorIs(err, ErrNegativeValue)
}

func (c *CopyDataSuite) TestErrorInvalidFile() {
	err := Copy("/dev/random", "", 0, 0)
	c.Require().ErrorIs(err, ErrUnsupportedFile)
}

func (c *CopyDataSuite) TestErrorCopyFile() {
	var pathErr *fs.PathError
	tmpName := "test"
	err := Copy(tmpName, tmpName, 0, 0)
	c.Require().ErrorAs(err, &pathErr)
}

func (c *CopyDataSuite) TestErrorLimit() {
	tmpName := generateFilename()
	err := Copy(c.f.Name(), tmpName, 0, -1)
	c.Require().ErrorIs(err, ErrNegativeValue)
}
