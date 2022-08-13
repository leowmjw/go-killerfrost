package teleport

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
