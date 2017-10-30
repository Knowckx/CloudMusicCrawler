package crawlerpkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

var modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
var nonce = "0CoJUm6Qyw8W8jud"
var pubKey = "010001"

//CreateSecretKey 创建加密秘钥
func CreateSecretKey(size int) string {
	keys := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := ""
	for i := 0; i < size; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano())) //生成随机数
		pos := r.Intn(size)                                  //产生一个0-size之间的数
		key += string(keys[pos])
	}
	return key
}

//PKCS5Padding padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	bufLen := len(ciphertext)
	padLen := blockSize - (bufLen % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(ciphertext, padText...)
}

//PKCS5UnPadding aes unpadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AesEncrypt aes-128 cbc iv == "0102030405060708"
func AesEncrypt(key []byte, text string) string {
	block, _ := aes.NewCipher(key)
	plianText := []byte(text)
	pad := PKCS5Padding(plianText, block.BlockSize())
	ciphertext := make([]byte, len(pad))

	cbc := cipher.NewCBCEncrypter(block, []byte("0102030405060708"))
	cbc.CryptBlocks(ciphertext, pad)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

//RsaEncrypt 对密钥进行加密处理
func RsaEncrypt(text, pubKey, modulus string) string {
	text = reverse(text)

	bigText := bigInt(fmt.Sprintf("%x", text))
	bigPk := bigInt(pubKey)
	bigMod := bigInt(modulus)

	b := bigText.Exp(bigText, bigPk, bigMod)
	return zfill(fmt.Sprintf("%x", b), 256)
}

// string reverse
func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
func bigInt(s string) *big.Int {
	big := new(big.Int)
	big.SetString(s, 16)
	return big
}
func zfill(str string, size int) string {
	for len(str) < size {
		str = "0" + str
	}
	return str
}

//CryptBody 加密的请求body
type CryptBody struct {
	Params    string `json:"params"`
	EncSecKey string `json:"encSecKey"`
}

//Encrypt 把post的结构体作为名文，进行加密参数
func Encrypt(params *RequestParam) (*CryptBody, error) {
	jsonByte, err := json.Marshal(params) //先把post参数改成json
	if err != nil {
		return nil, fmt.Errorf("Encrypt json.Marshal  Error:%v", err)
	}
	jsonStr := string(jsonByte) //换算成json字串

	secKey := CreateSecretKey(16) //创建一个16长的随机字串  密钥

	//1. Json数据第一次AES加密。密码为使用nonce。
	//2.把结果用随机生存的Key再AES一次
	cryptParams := AesEncrypt([]byte(secKey), AesEncrypt([]byte(nonce), jsonStr))
	//3.把AES使用的key使用Rsa加密。
	encSecKey := RsaEncrypt(secKey, pubKey, modulus)

	//新建结构体，把加密后的内容都放进去
	body := &CryptBody{Params: cryptParams, EncSecKey: encSecKey}
	return body, nil

}
