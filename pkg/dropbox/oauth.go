package dropbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joooostb/wahoo-garmin-sync/pkg/redis"
	log "github.com/sirupsen/logrus"
)

var API_URL string = "https://api.dropboxapi.com"
var APP_KEY = os.Getenv("APP_KEY")
var APP_SECRET = os.Getenv("APP_SECRET")
var APP_URL = os.Getenv("APP_URL")

func Authorize(c *gin.Context) error {
	c.Redirect(http.StatusTemporaryRedirect,
		fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?client_id=%s&redirect_uri=%s/oauth2&response_type=code&token_access_type=offline", APP_KEY, APP_URL))
	return nil
}

func Authenticate(c *gin.Context) error {
	if v, b := c.GetQuery("code"); b {
		_, err := authenticate(v, false)
		if err != nil {
			log.Error(err)
		}
		c.Status(http.StatusOK)
	}
	return nil
}

func Refresh(refreshToken string) (*Token, error) {
	t, err := authenticate(refreshToken, true)
	if err != nil {
		log.Error(err)
	}
	return t, err
}

func authenticate(token string, refresh bool) (*Token, error) {
	t := Token{}
	client := &http.Client{}

	var uri string
	if refresh {
		uri = fmt.Sprintf("code=%s&grant_type=refresh_token", token)
	} else {
		uri = fmt.Sprintf("code=%s&grant_type=authorization_code&redirect_uri=%s/oauth2", token, APP_URL)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/token", API_URL),
		strings.NewReader(uri))
	if err != nil {
		return nil, err
	}

	// Set base64 encoded Authorization header with APP_KEY & APP_SECRET
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", APP_KEY, APP_SECRET)))))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r, ctx, _ := redis.Default()
	if refresh {
		r.HSet(ctx, t.AccountID, "access_token", t.AccessToken, "expiry", t.ExpiresIn+time.Now().Unix())
	} else {
		r.HSet(ctx, t.AccountID, "access_token", t.AccessToken, "refresh_token", t.RefreshToken, "expiry", t.ExpiresIn+time.Now().Unix())
	}

	return &t, nil
}
