package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wms-server/constants"
	"wms-server/helpers/models"

	"github.com/sirupsen/logrus"
)

// ListErrorCode ...
var (
	ListErrorCode []models.MappingErrorCodes
)

// RegisterErrorCode registering code
func RegisterErrorCode() bool {

	var b []byte
	var err error
	b, err = os.ReadFile("errorcodes.json") // just pass the file name
	if err != nil {
		b, err = os.ReadFile("../errorcodes.json") // just pass the file name
		if err != nil {
			root, err := FindRoot()
			if err != nil {
				fmt.Println("Failed to find root folder : ", err)
			}
			b, err = os.ReadFile(root + "/errorcodes.json")
			if err != nil {
				fmt.Println("Failed to read file error code json : ", err)
			}
		}
	}

	if json.Unmarshal(b, &ListErrorCode) != nil {
		fmt.Println("Failed to unmarshaling json response error mapping ")
		return false
	}

	return true
}

// FindRoot ...
func FindRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", fmt.Errorf("could not find project root")
		}
		currentDir = parentDir
	}
}

// GetMetaResponse ..
func GetMetaResponse(key string) models.MetaData {
	var meta models.MetaData

	onceRegisterCode.Do(func() {
		RegisterErrorCode()
	})

	if key == constants.RC_SUCCESS {
		meta.Code = fmt.Sprintf("200%v00", constants.SERVICE_CODE)
		meta.Title = "Success"
		meta.Message = "Successful"
		return meta
	}

	for _, element := range ListErrorCode {
		if element.Key == key {
			meta.Code = strings.Replace(element.Content.Code, "?", constants.SERVICE_CODE, 1)
			meta.Title = element.Content.Title
			meta.Message = element.Content.Message
			return meta
		}
	}

	meta.Code = fmt.Sprintf("%v%v%v", http.StatusInternalServerError, constants.SERVICE_CODE, "00")
	meta.Title = "Error"
	meta.Message = "General error"
	return meta
}

// MapStatusCode ...
func MapStatusCode(metaCode models.MetaData) int {
	var output int
	switch {
	case strings.HasPrefix(metaCode.Code, "200"):
		output = http.StatusOK
	case strings.HasPrefix(metaCode.Code, "202"):
		output = http.StatusAccepted
	case strings.HasPrefix(metaCode.Code, "400"):
		output = http.StatusBadRequest
	case strings.HasPrefix(metaCode.Code, "401"):
		output = http.StatusUnauthorized
	case strings.HasPrefix(metaCode.Code, "403"):
		output = http.StatusForbidden
	case strings.HasPrefix(metaCode.Code, "404"):
		output = http.StatusNotFound
	case strings.HasPrefix(metaCode.Code, "409"):
		output = http.StatusConflict
	default:
		output = http.StatusInternalServerError
	}
	return output
}

// NEW //

// Struktur untuk menyimpan data dari file JSON
var metaData map[string]map[string]models.MetaData

// LoadMetaFiles memuat semua file JSON yang sesuai dengan pola errorcodes-*.json
func LoadMetaFiles(directory string) error {
	metaData = make(map[string]map[string]models.MetaData)

	files, err := filepath.Glob(filepath.Join(directory, "errorcodes-*.json"))
	if err != nil {
		logrus.Error("Error reading errorcodes files: ", err)
		return err
	}

	for _, file := range files {
		lang := extractLanguageCode(file)
		if lang == "" {
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			logrus.Errorf("Error reading file %s: %v", file, err)
			continue
		}

		var tempData []models.MappingErrorCodes
		err = json.Unmarshal(data, &tempData)
		if err != nil {
			logrus.Errorf("Error parsing JSON in file %s: %v", file, err)
			continue
		}

		// Simpan data ke dalam map dengan key sebagai indeks
		metaMap := make(map[string]models.MetaData)
		for _, entry := range tempData {
			metaMap[entry.Key] = entry.Content
		}
		metaData[lang] = metaMap
	}

	return nil
}

// extractLanguageCode mendapatkan kode bahasa dari nama file
func extractLanguageCode(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "-")
	if len(parts) < 2 {
		return ""
	}
	lang := strings.TrimSuffix(parts[1], filepath.Ext(parts[1]))
	return lang
}

// GetMetaResponse mengambil response berdasarkan bahasa, key, dan kode service
func GetNewMetaResponse(lang, key string) (response models.MetaData) {
	onceRegisterCode.Do(func() {
		LoadMetaFiles("./errorcodes")
	})
	serviceCode := constants.SERVICE_CODE
	if langData, exists := metaData[lang]; exists {
		if response, found := langData[key]; found {
			// Gantikan `?` dalam code dengan kode service
			response.Code = strings.Replace(response.Code, "?", serviceCode, 1)
			response.Message = fmt.Sprintf("%s (%s)", response.Message, response.Code)
			return response
		}
	}

	response.Code = "500" + serviceCode + "00"
	response.Title = "Error"
	response.Message = "General Error (500" + serviceCode + "00)"

	return response
}
