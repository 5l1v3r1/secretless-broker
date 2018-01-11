package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/kgilpin/secretless/internal/app/secretless/conjur"
)

var AdminAPIKey = os.Getenv("CONJUR_AUTHN_API_KEY")
var Host = os.Getenv("TEST_PROXY_HOST")

func psql(host string, port int, user string, environment []string) (string, error) {
	if host == "" {
		if Host != "" {
			host = Host
		} else {
			host = "secretless_test"
		}
	}

	args := []string{ "-h", host }
	if port != 0 {
		args = append(args, "-p")
		args = append(args, fmt.Sprintf("%d", port))
	}
	if user != "" {
		args = append(args, "-U")
		args = append(args, user)
	}
	args = append(args, "-c")
	args = append(args, "select count(*) from conjur.test")
	args = append(args, "sslmode=disable dbname=postgres")

	log.Println(strings.Join(append([]string{"psql"}, args...), " "))

	cmd := exec.Command("psql", args...)
	env := os.Environ()
	for _, v := range environment {
		env = append(env, v)
	}
	cmd.Env = env
	cmdOut, err := cmd.CombinedOutput()
	return string(cmdOut), err
}

func TestUnixSocketPasswordLogin(t *testing.T) {
	log.Print("Connect via Unix socket without authentication")

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmdOut, err := psql(fmt.Sprintf("%s/../../run/postgresql", cwd), 0, "", []string{})

	if err != nil {
		t.Fatal(cmdOut)
	}

	if !strings.Contains(cmdOut, "1 row") {
		t.Fatalf("Expected to find '1 row' in : %s", cmdOut)
	}
}

func TestStaticPasswordLogin(t *testing.T) {
	log.Print("Use a statically configured password")

	cmdOut, err := psql("", 5433, "alice", []string{"PGPASSWORD=alice"})

	if err != nil {
		t.Fatal(cmdOut)
	}

	if !strings.Contains(cmdOut, "1 row") {
		t.Fatalf("Expected to find '1 row' in : %s", cmdOut)
	}
}

func TestStaticPasswordLoginFailed(t *testing.T) {
	log.Print("Provide the wrong value for a statically configured password")

	cmdOut, err := psql("", 5433, "alice", []string{"PGPASSWORD=foobar"})

	if err == nil {
		t.Fatalf("Expected failed login : %s", cmdOut)
	}

	if !strings.Contains(cmdOut, "FATAL") {
		t.Fatalf("Expected to find 'FATAL' in : %s", cmdOut)
	}
	if !strings.Contains(cmdOut, "Login failed") {
		t.Fatalf("Expected to find 'Login failed' in : %s", cmdOut)
	}
}

func TestConjurLogin(t *testing.T) {
	log.Print("Provide a Conjur access token as the password")

	if AdminAPIKey == "" {
		t.Fatalf("CONJUR_AUTHN_API_KEY is missing")
	}

	var (
		adminToken *conjur.AccessToken
		userAPIKey string
		userToken  *conjur.AccessToken
		err        error
	)

	if adminToken, err = conjur.Authenticate("admin", AdminAPIKey); err != nil {
		t.Fatalf("Failed to authenticate as 'admin' : %s", err)
	}
	if userAPIKey, err = conjur.RotateAPIKey("bob", *adminToken); err != nil {
		t.Fatalf("Failed to rotate API key of user 'bob' : %s", err)
	}
	if userToken, err = conjur.Authenticate("bob", userAPIKey); err != nil {
		t.Fatalf("Failed to authenticate as 'bob' : %s", err)
	}

	userToken64 := base64.StdEncoding.EncodeToString([]byte(userToken.Token))

	cmdOut, err := psql("", 0, "bob", []string{fmt.Sprintf("PGPASSWORD=%s", userToken64)})

	if err != nil {
		t.Fatal(cmdOut)
	}

	if !strings.Contains(cmdOut, "1 row") {
		t.Fatalf("Expected to find '1 row' in : %s", cmdOut)
	}
}

func TestConjurUnauthorized(t *testing.T) {
	log.Print("Provide a Conjur access token for an unauthorized user")

	if AdminAPIKey == "" {
		t.Fatalf("CONJUR_AUTHN_API_KEY is missing")
	}

	var (
		adminToken *conjur.AccessToken
		userAPIKey string
		userToken  *conjur.AccessToken
		err        error
	)

	if adminToken, err = conjur.Authenticate("admin", AdminAPIKey); err != nil {
		t.Fatalf("Failed to authenticate as 'admin' : %s", err)
	}
	if userAPIKey, err = conjur.RotateAPIKey("charles", *adminToken); err != nil {
		t.Fatalf("Failed to rotate API key of user 'charles' : %s", err)
	}
	if userToken, err = conjur.Authenticate("charles", userAPIKey); err != nil {
		t.Fatalf("Failed to authenticate as 'charles' : %s", err)
	}

	userToken64 := base64.StdEncoding.EncodeToString([]byte(userToken.Token))

	cmdOut, err := psql("", 5434, "charles", []string{fmt.Sprintf("PGPASSWORD=%s", userToken64)})

	if err == nil {
		t.Fatal(cmdOut)
	}

	if !strings.Contains(cmdOut, "FATAL") {
		t.Fatalf("Expected to find 'FATAL' in : %s", cmdOut)
	}
	if !strings.Contains(cmdOut, "Conjur authorization failed") {
		t.Fatalf("Expected to find 'Conjur authorization failed' in : %s", cmdOut)
	}
}

func TestConjurLocal(t *testing.T) {
	log.Print("Proxy requires no authorization and will obtain its own Conjur access token")

	var (
		err error
	)

	cmdOut, err := psql("", 0, "", []string{})

	if err != nil {
		t.Fatal(cmdOut)
	}

	if !strings.Contains(cmdOut, "1 row") {
		t.Fatalf("Expected to find '1 row' in : %s", cmdOut)
	}
}