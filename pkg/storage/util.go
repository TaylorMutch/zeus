package storage

import (
	"fmt"
	"os"
)

func ReadObjstoreConfig(objstoreString, objstoreConfigFile string) ([]byte, error) {
	var objstoreConfig []byte
	if len(objstoreString) > 0 {
		objstoreConfig = []byte(objstoreString)
	} else {
		var err error
		objstoreConfig, err = os.ReadFile(objstoreConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read objstore config file: %w", err)
		}
	}

	return objstoreConfig, nil
}
