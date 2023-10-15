package encodeutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

func HexEncode(body []byte) string {
	return hex.EncodeToString(body)
}

func SNAPEncodingHex(body string) string {
	return strings.ToLower(HexEncode(Sha256Encode([]byte(MinifyJson(body)))))
}

func SNAPEncodingBase64(body string) string {
	return strings.ToLower(Base64(Sha256Encode([]byte(MinifyJson(body)))))
}

func Base64(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}

func Sha256Encode(body []byte) []byte {
	h := sha256.New()
	h.Write(body)
	return []byte(h.Sum(nil))
}

func MinifyJson(body string) string {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	body, err := m.String("application/json", body)
	if err != nil {
		panic(err)
	}

	return body
}

func HmacSha512(key []byte, data []byte) string {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
