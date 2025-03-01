package hw10programoptimization

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func (u *User) parse(s string) error {
	p := parserPool.Get()
	defer parserPool.Put(p)
	v, err := p.Parse(s)
	if err != nil {
		return err
	}

	u.ID = v.GetInt("Id")
	u.Name = string(v.GetStringBytes("Name"))
	u.Username = string(v.GetStringBytes("Username"))
	u.Email = string(v.GetStringBytes("Email"))
	u.Phone = string(v.GetStringBytes("Phone"))
	u.Password = string(v.GetStringBytes("Password"))
	u.Address = string(v.GetStringBytes("Address"))

	return nil
}
