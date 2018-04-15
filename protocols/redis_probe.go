//
// This is our redis protocol-test.
//
//
package protocols

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/skx/overseer/test"
)

//
// Our structure.
//
type REDISTest struct {
}

//
// Make a Redis-test against the given target.
//
func (s *REDISTest) RunTest(tst test.Test, target string, opts TestOptions) error {

	//
	// Predeclare our error
	//
	var err error

	//
	// The default port to connect to.
	//
	port := 6379

	//
	// The default password to use.
	//
	password := ""

	//
	// If the user specified a different port update to use it.
	//
	if tst.Arguments["port"] != "" {
		port, err = strconv.Atoi(tst.Arguments["port"])
		if err != nil {
			return err
		}
	}

	//
	// If the user specified a password use it.
	//
	password = tst.Arguments["password"]

	//
	// Default to connecting to an IPv4-address
	//
	address := fmt.Sprintf("%s:%d", target, port)

	//
	// If we find a ":" we know it is an IPv6 address though
	//
	if strings.Contains(target, ":") {
		address = fmt.Sprintf("[%s]:%d", target, port)
	}

	//
	// Attempt to connect to the host with the optional password
	//
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})

	//
	// And run a ping
	//
	// If the connection is refused, or the auth-details don't match
	// then we'll see that here.
	//
	_, err = client.Ping().Result()
	if err != nil {
		return err
	}

	//
	// If we reached here all is OK
	//
	return nil
}

//
// Register our protocol-tester.
//
func init() {
	Register("redis", func() ProtocolTest {
		return &REDISTest{}
	})
}