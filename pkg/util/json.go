package util

import (
	"encoding/json"
	"github.com/funnyecho/code-push/pkg/fs"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func JsonStringifyToFile(path string, data interface{}) error {
	plainData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "failed to json.stringify data: %v", data)
	} else {
		file, fileErr := fs.File(fs.FileConfig{FilePath: path})

		if fileErr != nil {
			return errors.Wrapf(fileErr, "invalid file, path:%s", path)
		}

		err := ioutil.WriteFile(
			file.Path(),
			plainData,
			os.ModePerm,
		)

		if err != nil {
			return errors.Wrapf(err, "failed to write json to path, path:%s, plainData:%s", path, plainData)
		}
	}

	return nil
}

func JsonParseFromFile(path string, dist interface{}) error {
	file, fileErr := fs.File(fs.FileConfig{FilePath: path})

	if fileErr != nil {
		return errors.Wrapf(fileErr, "invalid file, path:%s", path)
	}

	plainIndex, err := ioutil.ReadFile(file.Path())
	if err == nil {
		err = json.Unmarshal(plainIndex, dist)
		if err != nil {
			return errors.Wrapf(err, "failed to json.parse from file, path:%s, plainIndex:%s", path, plainIndex)
		}
	} else {
		return errors.Wrapf(err, "failed to load json file, path:%", path)
	}

	return nil
}

func JsonStringify(data interface{}) ([]byte, error) {
	plainData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json.stringify data, data:%v", data)
	}

	return plainData, nil
}

func JsonParse(data []byte, dist interface{}) error {
	err := json.Unmarshal(data, dist)
	if err != nil {
		return errors.Wrapf(err, "failed to json.parse, data:%s", data)
	}

	return nil
}
