package middlewares

import (
	"github.com/ChristopherRabotin/gin-contrib-headerauth"
	"errors"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"
)

type SHA384Manager struct {
	// --> If using a database to check for the secret, you'll probably use a different struct, which may have a pointer
	// --> to your database connection or even not set it, and have all the database connection, querying, and friends
	// --> performed in the `CheckHeader` function.
	Secret string
	*headerauth.HMACManager
}

// CheckHeader returns the secret key and the data to sign from the provided access key.
// Here should reside additional verifications on the header, or other parts of the request, if needed.
func (m SHA384Manager) CheckHeader(auth *headerauth.AuthInfo, req *http.Request) (err *headerauth.AuthErr) {
	if req.ContentLength != 0 && req.Body == nil {
		// Not sure whether net/http or Gin handles these kinds of fun situations.
		return &headerauth.AuthErr{400, errors.New("received a forged packet")}
	}
	// Grabbing the date and making sure it's in the correct format and is within fifteen minutes.
	dateHeader := req.Header.Get("RequestDate")
	if dateHeader == "" {
		return &headerauth.AuthErr{406, errors.New("no Date header provided")}
	}
	date, derr := time.Parse("2006-01-02T15:04:05.000Z", dateHeader)
	if derr != nil {
		return &headerauth.AuthErr{408, errors.New("could not parse date")}
	} else if time.Since(date) > time.Minute*15 {
		return &headerauth.AuthErr{410, errors.New("request is too old")}
	}

	// --> Here is where you would do a database call to check if the access key is valid
	// --> and what the appropriate secret key is, e.g.:
	// if secretKey, dbErr := getSecretFromDB(access); dbErr == nil && auth.Secret == secretKey { ...
	if auth.AccessKey == "my_access_key" {
		// In this example, we'll be implementing a *similar* signing method to the Amazon AWS REST one.
		// We'll use the HTTP-Verb, the MD5 checksum of the Body, if any, and the Date header in ISO format.
		// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
		// Note: We are returning a variety of error codes which don't follow the spec only for the purpose of testing.
		serializedData := req.Method + "\n"
		if req.ContentLength != 0 {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return &headerauth.AuthErr{402, errors.New("could not read the body")}
			}
			hash := md5.New()
			hash.Write(body)
			serializedData += hex.EncodeToString(hash.Sum(nil)) + "\n"
		} else {
			serializedData += "\n"
		}
		// We know from Authorize that the Date header is present and fits our time constaints.
		serializedData += req.Header.Get("Date")

		auth.Secret = m.Secret
		auth.DataToSign = serializedData
		return
	}
	return &headerauth.AuthErr{418, errors.New("you are a teapot")}
}

func (m SHA384Manager) Authorize(auth *headerauth.AuthInfo) (interface{}, *headerauth.AuthErr) {
	if auth.AccessKey == "my_access_key" {
		return "All good with my access key!", nil
	}
	return "All good with any access key!", nil
}



// TokenManger is an example definition of an AuthKeyManager struct.
type TokenManger struct {
	*headerauth.TokenManager
}

// Authorize returns the secret key from the provided access key.
func (m TokenManger) CheckHeader(auth *headerauth.AuthInfo, req *http.Request) (err *headerauth.AuthErr) {
	auth.Secret = ""     // There is no secret key, just an access key.
	auth.DataToSign = "" // There is no data to sign in Token auth.
	if auth.AccessKey != "valid" {
		err = &headerauth.AuthErr{403, errors.New("invalid access key")}
	}
	return
}

func (m TokenManger) Authorize(auth *headerauth.AuthInfo) (val interface{}, err *headerauth.AuthErr) {
	return true, nil
}




// HTTPBasic is an example of an HTTP Basic Auth.
type HTTPBasicManager struct {
	Accounts map[string]string // --> Here we are using a hard coded map, but the logic is up to the dev.
	*headerauth.HTTPBasicAuth // Embedded struct greatly helps in defining HTTP Basic Auth.
}

// Authorize checks that the provided authorization is valid.
// --> Here is where you can interface with a database, or something which stores the list of valid usernames
// --> and their associated passwords. Note that in the other schemes we try to fail earlier (in CheckHeader).
func (m HTTPBasicManager) Authorize(auth *headerauth.AuthInfo) (val interface{}, err *headerauth.AuthErr) {
	if password, ok := m.Accounts[auth.AccessKey]; !ok || password != auth.Secret {
		err = &headerauth.AuthErr{401, errors.New("invalid credentials")}
	} else {
		// In CheckHeader we changed the AccessKey to be the actual username, instead
		// of the Base64 encoded authentication string.
		val = auth.AccessKey
	}
	return
}