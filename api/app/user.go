package app

import (
	"fmt"
	"strings"

	"github.com/okteto/app/api/k8s/client"
	"github.com/okteto/app/api/k8s/serviceaccounts"
	"github.com/okteto/app/api/model"
)

//CreateUser configures a service account for a given user
func CreateUser(u *model.User) error {
	c, err := client.Get()
	if err != nil {
		return fmt.Errorf("error getting k8s client: %s", err)
	}

	if err := serviceaccounts.Create(u, c); err != nil {
		return err
	}

	s := model.NewSpace(u.GithubID, u, []model.Member{})
	if err := CreateSpace(s); err != nil {
		return err
	}
	return nil
}

//GetCredential returns the credentials of the user for her space
func GetCredential(u *model.User, space string) (string, error) {
	credential, err := serviceaccounts.GetCredentialConfig(u, space)
	if err != nil {
		return "", err
	}

	return credential, err
}

//FindOrKeepUser retrieves user if it exists
func FindOrKeepUser(u *model.User) (*model.User, error) {
	found, err := serviceaccounts.GetUserByGithubID(u.GithubID)
	if err == nil {
		return found, nil
	}
	if !strings.Contains(err.Error(), "not found") {
		return nil, err
	}
	return u, nil
}
