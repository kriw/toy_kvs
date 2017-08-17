package util

import (
	"bytes"
	"github.com/ncw/directio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var isDirect = true

func align(n int) (aligned int) {
	n >>= 1
	aligned = 1
	for n > 0 {
		aligned <<= 1
		n >>= 1
	}
	return
}

func WriteFile(fileName string, content []byte) {
	if isDirect {
		out, err := directio.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
		defer out.Close()
		if err != nil {
			LogFatal("Failed to directio.OpenFile for read", err)
		}
		pad := make([]byte, 4096-len(content)%4096)
		_, err = out.Write(append(content, pad...))
		if err != nil {
			LogFatal("Failed to write", err)
		}
	} else {
		ioutil.WriteFile(fileName, content, os.ModePerm)
	}
}

func ReadFile(fileName string) ([]byte, error) {
	if isDirect {
		// Read the file
		in, err := directio.OpenFile(fileName, os.O_RDONLY, 0666)
		defer in.Close()
		if err != nil {
			LogFatal("Failed to directio.OpenFile for write", err)
		}
		var size int64
		if fi, err := in.Stat(); err == nil {
			size = fi.Size()
		}
		blockSize := 1
		tmp := size
		for tmp > 0 {
			blockSize <<= 1
			tmp >>= 1
		}
		block := directio.AlignedBlock(blockSize)
		_, err = io.ReadAtLeast(in, block, int(size))
		return block[:size+bytes.MinRead], err
	} else {
		return ioutil.ReadFile(fileName)
	}
}

func ReadDir(dirName string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirName)
}

func MatchExt(filename string, ext string) bool {
	ns := strings.Split(filename, ".")
	return ns[len(ns)-1] == ext
}

func FilesMap(dirName string, f func(string)) {
	fileList, _ := ioutil.ReadDir(dirName)
	for _, file := range fileList {
		fileName := file.Name()
		f(fileName)
	}
}
