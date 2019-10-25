package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"

	"github.com/ctailor2/clirescue/cmdutil"
	"github.com/ctailor2/clirescue/user"
)

var (
	url          string     = "https://www.pivotaltracker.com/services/v5/me"
	fileLocation string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	stdout       *os.File   = os.Stdout
)

// Me - does stuff
func Me() {
	setCredentials()
	parse(makeRequest())
	ioutil.WriteFile(fileLocation, []byte(currentUser.APIToken), 0644)
}

func makeRequest() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(currentUser.Username, currentUser.Password)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) {
	var meResp = new(meResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	currentUser.APIToken = meResp.APIToken
}

func setCredentials() {
	fmt.Fprint(stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type meResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}
