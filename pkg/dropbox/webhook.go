package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	connect "github.com/abrander/garmin-connect"
	"github.com/gin-gonic/gin"
	"github.com/joooostb/wahoo-garmin-sync/pkg/redis"
	log "github.com/sirupsen/logrus"
)

// Respond to the webhook verification (GET request)
// by echoing back the challenge parameter.
func Challenge(c *gin.Context) error {

	if v, b := c.GetQuery("challenge"); b {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, "text/plain", []byte(v))
	}
	return nil
}

func Handler(c *gin.Context) error {
	var n Notification
	if c.Request.Method == http.MethodPost {
		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &n)
		log.Debug(string(body))
		if err != nil {
			log.Errorf("failed to unmarshal Dropbox info: %s", string(body))
		}

		fit, err := getFit(n.ListFolder.Accounts[0])
		if err != nil {
			log.Error(err)
			return err
		}

		c := connect.NewClient()
		c.Email = "REDACTED"
		c.Password = "REDACTED"
		err = c.Authenticate()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debug("uploading file to Garmin Connect")
		c.ImportActivity(bytes.NewReader(fit), connect.ActivityFormatFIT)
	}
	return nil
}

// Get Access Token for account ID from Redis
func getToken(id string) (*Token, error) {
	r, ctx, _ := redis.Default()
	// Use cached token if not expired
	if r.HGet(ctx, id, "expiry").Val() > fmt.Sprint(time.Now().Unix()) {
		return &Token{
			AccessToken:  r.HGet(ctx, id, "access_token").Val(),
			RefreshToken: r.HGet(ctx, id, "refresh_token").Val(),
		}, nil
	} else {
		t, _ := Refresh(r.HGet(ctx, id, "refresh_token").Val())
		return t, nil
	}

}

// id: UserID
// c: Continue
func listFolder(id string) (ListFolderResponse, error) {
	l := ListFolderResponse{}

	token, err := getToken(id)
	if err != nil {
		log.Errorf("failed to get token for user: %s", id)
	}

	client := &http.Client{}
	r, ctx, _ := redis.Default()

	has_more := true
	cursor := ""
	for has_more {
		// If cursor is set, change URL
		uri := fmt.Sprintf("%s/2/files/list_folder", API_URL)
		v, _ := json.Marshal(ListFolderRequest{
			Path: "/Apps/WahooFitness",
		})
		if cursor != "" {
			uri = fmt.Sprintf("%s/2/files/list_folder/continue", API_URL)
			cursor = r.HGet(ctx, id, "cursor").Val()
			v, _ = json.Marshal(ListFolderRequestContinue{
				Cursor: cursor,
			})
		}
		req, err := http.NewRequest("POST", uri, strings.NewReader(string(v)))

		if err != nil {
			log.Error(err)
			return l, err
		}
		log.Debugf("Token used: %s", token)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return l, err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		log.Debug(string(body))
		if err != nil {
			fmt.Println(err)
			return l, err
		}

		m := ListFolderResponse{}
		err = json.Unmarshal(body, &m)
		has_more = m.HasMore

		if err != nil {
			fmt.Println(err)
			return l, err
		}
		l.Entries = append(l.Entries, m.Entries...)
	}
	return l, err
}

func getFit(id string) ([]byte, error) {
	file := []byte{}
	l, err := listFolder(id)
	if err != nil {
		log.Error(err)
		return file, err
	}

	sort.Slice(l.Entries, func(i, j int) bool {
		return l.Entries[i].ClientModified.Before(l.Entries[j].ClientModified)
	})

	last := l.Entries[len(l.Entries)-1]
	r, ctx, _ := redis.Default()

	// If most recent hash is the same as last time, do nothing
	if last.ContentHash == r.HGet(ctx, id, "last_hash").Val() {
		log.Debug("skipping as file with this hash has already been processed: %s / %s", last.Name, last.ContentHash)
		return file, nil
	} else {
		log.Debugf("found latest .fit file: %s", last.Name)
		r.HSet(ctx, id, "last_hash", last.ContentHash)
	}
	// Something here to download the file from Dropbox
	uri := fmt.Sprintf("https://content.dropboxapi.com/2/files/download")

	v := fmt.Sprintf("{\"path\": \"%s\"}", last.ID)
	req, err := http.NewRequest("POST", uri, io.Reader(nil))

	if err != nil {
		log.Error(err)
		return file, err
	}

	token, _ := getToken(id)
	log.Debugf("Token used: %s", token)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Dropbox-API-Arg", string(v))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return file, err
	}

	defer res.Body.Close()
	file, err = ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return file, err
	}

	return file, nil
}
