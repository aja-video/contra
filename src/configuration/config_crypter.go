package configuration

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"log"
	"reflect"
	"strings"
)

// EncryptConfigFile will encrypt all the passwords in the full file.
func EncryptConfigFile(filePath string) {
	saveRequired := false

	log.Println("ENCRYPT: Checking file: " + filePath)
	iniFile, err := ini.Load(filePath)
	if err != nil {
		log.Println("ENCRYPT: Failed to load: " + filePath)
		panic(err) // Critical....
		//return nil, err
	}

	// Find or create a config key.
	if !iniFile.Section("main").HasKey("EncryptKey") {
		// Generate new key.
		keyNew := make([]byte, 32)
		_, err := rand.Read(keyNew)
		if err != nil {
			panic(err) // TODO...
		}

		// Store it
		iniFile.Section("main").NewKey("EncryptKey", fmt.Sprintf("%x", keyNew))
		saveRequired = true
	}

	// If they muck with their key, there will be problems in the encryption step.
	// Not much we can do...
	encryptKey := iniFile.Section("main").Key("EncryptKey").Value()

	// Core assumption, we will encrypt anything that the key includes "Pass" like "SMTPPass"
	for _, section := range iniFile.Sections() {
		for _, key := range section.Keys() {
			if strings.Contains(key.Name(), "Pass") {
				if !strings.Contains(key.Value(), "~~contra~~") {
					log.Printf("ENCRYPT: Pass Key Found - S: %s K: %s\n", section.Name(), key.Name())
					key.SetValue("~~contra~~" + encryptConfig(encryptKey, key.Value()))
					saveRequired = true
				}
			}
		}
	}

	if saveRequired {
		// Save our changes to the config file.
		log.Println("ENCRYPT: Change detected, saving file: " + filePath)
		ini.PrettyFormat = true
		ini.PrettySection = true
		iniFile.SaveToIndent(filePath, "    ")
	}

	log.Println("ENCRYPT: Done with file: " + filePath)
}

func decryptLoadedConfig(config *Config) {

	fmt.Println(config)

	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		// Check it!
		val := fmt.Sprintf("%s", v.Field(i).Interface())
		if strings.HasPrefix(val, "~~contra~~") {
			// Need to decide...
			val = strings.Replace(val, "~~contra~~", "", 1)
			v.Field(i).SetString(decryptConfig(config.EncryptKey, val))
		}
	}

	// And the device configs...
	for id, device := range config.Devices {
		v := reflect.ValueOf(&device).Elem()
		for i := 0; i < v.NumField(); i++ {
			// Check it!
			val := fmt.Sprintf("%s", v.Field(i).Interface())
			if strings.HasPrefix(val, "~~contra~~") {
				// Need to decode...
				val = strings.Replace(val, "~~contra~~", "", 1)
				v.Field(i).SetString(decryptConfig(config.EncryptKey, val))
			}
		}

		// Seems like brute forcing reflection in a bad way, but this works!
		config.Devices[id] = v.Interface().(DeviceConfig)
	}
}

func encryptConfig(key, value string) string {
	// This forces the byte array length to 32.
	// The extra 0 padding is fine for our needs.
	keyByte := make([]byte, 32)
	copy(keyByte, key)

	return encryptSimple(keyByte, value)
}

func decryptConfig(key, value string) string {
	// This forces the byte array length to 32.
	// The extra 0 padding is fine for our needs.
	keyByte := make([]byte, 32)
	copy(keyByte, key)

	return decryptSimple(keyByte, value)
}

// EncryptSimple string to base64 crypto using AES
// - Concept from: https://gist.github.com/manishtpatel/8222606
func encryptSimple(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// DecryptSimple from base64 to decrypted string
func decryptSimple(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
