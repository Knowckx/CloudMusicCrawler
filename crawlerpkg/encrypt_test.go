package crawlerpkg

import (
	"testing"
)

func Test_createSecretKey(t *testing.T) {
	size := 16
	key := CreateSecretKey(size)
	if len(key) == size {
		t.Log("Test_createSecretKey success", key)
	} else {
		t.Error("Test_createSecretKey error", key)
	}
}

func Test_aesEncrypt(t *testing.T) {
	text := "djfjiojoiwjoijfoijiji"
	key := CreateSecretKey(16)
	aesText := AesEncrypt(text, key)
	decryptText := AesDecrypt(aesText, key)
	if decryptText == text {
		t.Log("Success --- originText:", text, "cryptText:", aesText, "decryptText:", decryptText)
	} else {
		t.Error("fail--- originText:", text, "cryptText:", aesText, "decryptText:", decryptText)
	}
}

func Test_Reverse(t *testing.T) {
	re := Reverse("qwertyuiopqaz")
	if re == "zaqpoiuytrewq" {
		t.Log("ok", re)
	} else {
		t.Error("fail", re)
	}
}

func Test_RsaEncrypt(t *testing.T) {
	s := RsaEncrypt("jfajofjiojfoijiofjwiojiwfkjkl", "010001", "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7")
	t.Log(s)
}
