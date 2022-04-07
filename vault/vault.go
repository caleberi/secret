package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"secret/encrypt"
	"sync"
)

type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	fmutex      sync.Mutex
	keyValues   map[string]string
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dc := json.NewDecoder(r)
	return dc.Decode(&v.keyValues)
}

func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

func MemoryVault(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string, 5),
	}
}

func (v *Vault) GetKeyInMemoryVault(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret :no value for that key 🌶️")
	}
	decryptedVal, err := encrypt.Decrypt(v.encodingKey, hex)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return decryptedVal, nil
}

func (v *Vault) SetKeyInMemoryVault(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	encryptedVal, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedVal
	return nil
}

func FileVault(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

func (v *Vault) loadKeyValues() error {
	v.fmutex.Lock()
	defer v.fmutex.Unlock()
	f, err := os.Open(v.filepath)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		fmt.Println("file does not exist but we are creating it ...🍔")
	// 		f, err = os.OpenFile(v.filepath, os.O_CREATE, 0755)
	// 		if err != nil {
	// 			v.keyValues = make(map[string]string)
	// 		}
	// 	} else {
	// 		v.keyValues = make(map[string]string)
	// 		return nil
	// 	}
	// }
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := encrypt.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.readKeyValues(r)

}

func (v *Vault) GetKeyInFileVault(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret :no value for that key 🌶️")
	}
	return hex, nil
}

func (v *Vault) SetKeyInFileVault(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	return v.saveKeyValues()
}

func (v *Vault) saveKeyValues() error {
	v.fmutex.Lock()
	defer v.fmutex.Unlock()
	f, err := os.OpenFile(v.filepath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := encrypt.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.writeKeyValues(w)
}
