package docs

import (
	"embed"
	jsoniter "github.com/json-iterator/go"
)

////go:embed generated_doc.swagger.json
//var f embed.FS

//go:embed swagger.json
var f embed.FS

func init() {
	var (
		docTemplMap = make(map[string]interface{})
	)
	fileData, _ := f.ReadFile("swagger.json")
	_ = jsoniter.Unmarshal(fileData, &docTemplMap)
	for k, v := range docTemplMap["paths"].(map[string]interface{}) {
		docTemplMap["paths"].(map[string]interface{})[k] = v
	}

}

func GetAPIData() map[string]any {
	var (
		docTemplMap = make(map[string]interface{})
	)
	customizeDoc, _ := f.ReadFile("swagger.json")
	_ = jsoniter.Unmarshal(customizeDoc, &docTemplMap)

	for k, v := range docTemplMap["paths"].(map[string]interface{}) {
		docTemplMap["paths"].(map[string]interface{})[k] = v
	}
	return docTemplMap
}
