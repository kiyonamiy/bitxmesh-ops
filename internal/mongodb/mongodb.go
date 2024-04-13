package mongodb

type Connection struct {
	host     string
	port     string
	user     string
	password string
}

func NewConnection(host string, port string, user string, password string) *Connection {
	return &Connection{
		host,
		port,
		user,
		password,
	}
}

func (c *Connection) Host() string {
	return c.host
}

func (c *Connection) Port() string {
	return c.port
}

func (c *Connection) User() string {
	return c.user
}

func (c *Connection) Password() string {
	return c.password
}
