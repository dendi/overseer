package protocols

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/skx/overseer/test"
)

//
// Our structure.
//
type XMPPTest struct {
}

//
// Run the test against the specified target.
//
func (s *XMPPTest) RunTest(tst test.Test, target string, opts TestOptions) error {
	var err error

	//
	// The default port to connect to.
	//
	port := 5222

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
	// Set an explicit timeout
	//
	d := net.Dialer{Timeout: opts.Timeout}

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
	// Make the TCP connection.
	//
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return err
	}

	//
	// Send a (bogus) greeting
	//
	_, err = conn.Write([]byte("<>\n"))
	if err != nil {
		return err
	}

	//
	// Read the response.
	//
	banner, err := bufio.NewReader(conn).ReadString('>')
	if err != nil {
		return err
	}

	//
	// Now close the conneciton.
	//
	err = conn.Close()
	if err != nil {
		return err
	}

	if !strings.Contains(banner, "<?xml") {
		return errors.New(fmt.Sprintf("Banner doesn't look like an XMPP-banner '%s'", banner))
	}

	return nil
}

//
// Register our protocol-tester.
//
func init() {
	Register("xmpp", func() ProtocolTest {
		return &XMPPTest{}
	})
}