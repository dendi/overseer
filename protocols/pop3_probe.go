package protocols

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/simia-tech/go-pop3"
	"github.com/skx/overseer/test"
)

//
// Our structure.
//
type POP3Test struct {
}

//
// Run the test against the specified target.
//
func (s *POP3Test) RunTest(tst test.Test, target string, opts TestOptions) error {
	var err error

	fmt.Printf("target:%s test.target:%s\n", target, tst.Target)

	//
	// The default port to connect to.
	//
	port := 110

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
	// Connect
	//
	c, err := pop3.Dial(address, pop3.UseTimeout(opts.Timeout))
	if err != nil {
		return err
	}

	//
	// Did we get a username/password?  If so try to authenticate
	// with them
	//
	if (tst.Arguments["username"] != "") && (tst.Arguments["password"] != "") {
		err = c.Auth(tst.Arguments["username"], tst.Arguments["password"])
		if err != nil {
			return err
		}
	}

	//
	// Quit and return
	//
	c.Quit()
	return nil
}

//
// Register our protocol-tester.
//
func init() {
	Register("pop3", func() ProtocolTest {
		return &POP3Test{}
	})
}