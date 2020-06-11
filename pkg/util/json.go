package util

import (
	"encoding/json"
	"github.com/funnyecho/code-push/pkg/errors"
	"github.com/funnyecho/code-push/pkg/fs"
	"io/ioutil"
	"os"
)

func JsonStringifyToFile(path string, data interface{}) error {
	plainData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Throw(errors.CtorConfig{
			Error: err,
			Msg:   "failed to json.stringify data",
			Meta:  errors.MetaFields{"data": data},
		})
	} else {
		file, fileErr := fs.File(fs.FileConfig{FilePath: path})

		if fileErr != nil {
			return errors.Throw(errors.CtorConfig{
				Error: fileErr,
				Msg:   "invalid file",
				Meta:  errors.MetaFields{"path": path},
			})
		}

		err := ioutil.WriteFile(
			file.Path(),
			plainData,
			os.ModePerm,
		)

		if err != nil {
			return errors.Throw(errors.CtorConfig{
				Error: err,
				Msg:   "failed to write json to path",
				Meta:  errors.MetaFields{"path": file.Path(), "plainData": plainData},
			})
		}
	}

	return nil
}

func JsonParseFromFile(path string, dist interface{}) error {
	file, fileErr := fs.File(fs.FileConfig{FilePath: path})

	if fileErr != nil {
		return errors.Throw(errors.CtorConfig{
			Error: fileErr,
			Msg:   "invalid file",
			Meta:  errors.MetaFields{"path": path},
		})
	}

	plainIndex, err := ioutil.ReadFile(file.Path())
	if err == nil {
		err = json.Unmarshal(plainIndex, dist)
		if err != nil {
			return errors.Throw(errors.CtorConfig{
				Error: err,
				Msg:   "failed to json.parse from file",
				Meta:  errors.MetaFields{"plainIndex": plainIndex, "path": file.Path()},
			})
		}
	} else {
		return errors.Throw(errors.CtorConfig{
			Error: err,
			Msg:   "failed to load json file",
			Meta:  errors.MetaFields{"path": file.Path()},
		})
	}

	return nil
}

func JsonStringify(data interface{}) ([]byte, error) {
	plainData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.Throw(errors.CtorConfig{
			Error: err,
			Msg:   "failed to json.stringify data",
			Meta:  errors.MetaFields{"data": data},
		})
	}

	return plainData, nil
}

func JsonParse(data []byte, dist interface{}) error {
	err := json.Unmarshal(data, dist)
	if err != nil {
		return errors.Throw(errors.CtorConfig{
			Error: err,
			Msg:   "failed to json.parse",
			Meta:  errors.MetaFields{"data": data},
		})
	}

	return nil
}
