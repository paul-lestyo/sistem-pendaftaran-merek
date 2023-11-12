package helper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BrandPDKI struct {
	BrandName  string `json:"brand_name"`
	Owner      string `json:"owner"`
	Similarity string `json:"similarity"`
}

func GetDataSearchPDKI(search string) []BrandPDKI {
	url := "https://pdki-indonesia-api.dgip.go.id/api/trademark/search2?keyword=" + url.QueryEscape(search) + "&page=1&type=trademark&order_state=asc"
	currentTime := time.Now()
	formattedDate := currentTime.Format("20060102")

	inputBytes := []byte(formattedDate + search + "arhoBkmdhcHsWSJPyLhLVqGNhAEontUgedqsNAAWjRkXkKDnrnNwolYRnEwGkaYaC")
	hash := sha256.New()
	hash.Write(inputBytes)
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)

	var jsonStr = []byte(`{"key": "` + hashString + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json := string(body)

	var brandData []BrandPDKI

	for i := 0; i < 3; i++ {
		brand := BrandPDKI{
			BrandName: gjson.Get(json, fmt.Sprintf("hits.hits.%d._source.image.0.brand_name", i)).String(),
			Owner:     gjson.Get(json, fmt.Sprintf("hits.hits.%d._source.owner.0.tm_owner_name", i)).String(),
		}
		if brand.BrandName != "" && brand.Owner != "" {
			brand.Similarity = CheckSimilarity(search, brand.BrandName)
			brandData = append(brandData, brand)
		}
	}

	return brandData
}

func CheckSimilarity(brandRegister string, brandRegistered string) string {
	similarity := strutil.Similarity(
		strings.ToLower(brandRegister),
		strings.ToLower(brandRegistered),
		metrics.NewJaroWinkler(),
	)

	intValue := int(similarity * 100)
	stringSimilarity := strconv.Itoa(intValue)
	return stringSimilarity + "%"
}
