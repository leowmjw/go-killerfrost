package breakglass

import "fmt"

type Bastion struct {
	Identity string
	URL      string
	Token    string
}

func (b Bastion) AddToRole(roleName string) error {
	return nil
}

func (b Bastion) RemoveFromRole(roleName string) error {
	return nil
}

func (b Bastion) GetRole() (string, error) {
	var roleName string

	// Query Teleport role ...
	if roleName == "" {
		return "", fmt.Errorf("INVLID ROLE ..!! %v", b)
	}
	return roleName, nil
}
