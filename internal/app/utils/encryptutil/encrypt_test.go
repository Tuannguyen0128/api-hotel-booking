package encryptutil

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
)

func Test_custom_happy_case(t *testing.T) {
	u := util{}
	x := u.NewEncryptor("6DF23442B9CF99C628138DCD77526")
	fmt.Println(x.Encrypt("secret", false))

	enc := u.NewEncryptor("key")
	testString := "string"
	type Data struct {
		A string
		B string
	}
	testData := Data{
		A: "a",
		B: "B",
	}

	{
		deterministic1, pe := enc.Encrypt(testString, true)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}
		deterministic2, pe := enc.Encrypt(testString, true)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}

		if deterministic1 != deterministic2 {
			t.Errorf("encrypt deterministic string give differance result")
		}

		data, pe := enc.Decrypt(deterministic1)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}

		if data != testString {
			t.Errorf("decrypted string data is not the same as source")
		}
	}
	{
		nonDeterministic1, pe := enc.Encrypt(testString, false)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}
		nonDeterministic2, pe := enc.Encrypt(testString, false)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}

		if nonDeterministic1 == nonDeterministic2 {
			t.Errorf("encrypt non deterministic string give same result")
		}

		data1, pe := enc.Decrypt(nonDeterministic1)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}

		if data1 != testString {
			t.Errorf("decrypted string data is not the same as source")
		}

		data2, pe := enc.Decrypt(nonDeterministic1)
		if pe != nil {
			t.Errorf("cannot deterministic encrypt string")
		}

		if data2 != testString {
			t.Errorf("decrypted string data is not the same as source")
		}
	}
	{
		eData, pe := enc.EncryptWithMarshalling(testData, false)
		if pe != nil {
			t.Errorf("cannot encrypt object")
		}

		var data Data
		pe = enc.DecryptWithUnmarshalling(eData, &data)
		if pe != nil {
			t.Errorf("cannot dencrypt object")
		}

		if !reflect.DeepEqual(testData, data) {
			t.Errorf("decrypted object data is not the same as source")
		}
	}
}

func Test_custom_bad_case(t *testing.T) {
	u := util{}
	enc := u.NewEncryptor("key")

	{
		fn := func() {}
		_, err := enc.EncryptWithMarshalling(fn, false)
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}
	}

	{
		_, err := enc.Decrypt("_")
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}

		_, err = enc.Decrypt(base64.StdEncoding.EncodeToString([]byte("test")))
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}

		_, err = enc.Decrypt(base64.StdEncoding.EncodeToString([]byte("testtesttesttesttesttesttesttesttesttesttest")))
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}

		err = enc.DecryptWithUnmarshalling("_", nil)
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}

		err = enc.DecryptWithUnmarshalling(base64.StdEncoding.EncodeToString([]byte("test")), nil)
		if err == nil {
			t.Errorf("expect error DecryptCiphertextInvalid")
		}

		ct, _ := enc.Encrypt("test", false)

		err = enc.DecryptWithUnmarshalling(ct, nil)
		if err == nil {
			t.Errorf("expect error DecryptPlaintextUnmarshal")
		}
	}
}
