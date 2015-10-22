package endpoints

import (
	"bytes"
	"errors"
	"log"
	"strconv"

	"github.com/FoxComm/libs/logger"
	"github.com/mailgun/vulcan/location/httploc"
	"github.com/mailgun/vulcan/middleware"

	"fmt"
	"net"
	"net/http"
	"net/url"
)

var Endpoints = []*Endpoint{}

func Add(ep *Endpoint) {
	Endpoints = append(Endpoints, ep)
}

type Endpoint struct {
	Name        string
	Description string
	Domain      string
	DefaultPort string
	APIPrefix   string
	routePrefix string
	IsFeature   bool
	Options     httploc.Options
	MiddleWares []middleware.Middleware
	Observers   []middleware.Observer
}

type WrappedURL url.URL

type encoding int

const (
	encodePath encoding = 1 + iota
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

func (ep *Endpoint) Port() string {
	if u, err := url.Parse(ep.Domain); err == nil {
		if _, port, err := net.SplitHostPort(u.Host); err == nil {
			return port
		} else {
			return ep.DefaultPort
		}
	}
	return "80"
}

func (ep *Endpoint) AppServerUrl() string {
	//Let's simply increment the port for now
	var newEp Endpoint
	newEp = *ep

	if u, err := url.Parse(newEp.Domain); err == nil {
		if host, port, err := net.SplitHostPort(u.Host); err == nil {
			newEpPortInt, _ := strconv.Atoi(port)
			newEpPortInt += 1
			newEpPortStr := strconv.Itoa(newEpPortInt)
			newEpHost := net.JoinHostPort(host, newEpPortStr)
			return fmt.Sprintf("%v%v", "http://", newEpHost)
		}
	}

	return ""
}

func (ep *Endpoint) RoutePrefix() string {
	if ep.routePrefix != "" {
		return ep.routePrefix
	} else {
		return ep.APIPrefix
	}
}

func (ep *Endpoint) Url() string {
	url, err := url.Parse(fmt.Sprintf("%v%v%v", "http://", ep.Domain, ep.APIPrefix))
	if err == nil {
		return url.String()
	} else {
		logger.Error("Invalid URL for: " + ep.Name)
		return ""
	}
}

func (ep *Endpoint) UrlWithoutScheme() string {
	if url, err := url.Parse(fmt.Sprintf("%v%v%v", "http://", ep.Domain, ep.APIPrefix)); err == nil {
		return (*WrappedURL)(url).StringWithoutScheme()
	}
	return ""
}

func (ep *Endpoint) GetHostWithPort() string {
	return ep.Domain + ":" + ep.DefaultPort
}

func (ep *Endpoint) Host() string {
	if url, err := url.Parse(fmt.Sprintf("%v%v", ep.Domain, ep.APIPrefix)); err == nil {
		return url.Host
	} else {
		log.Printf("EndpointHostError=%s", err)
	}
	return ""
}

func (ep *Endpoint) UseMiddleWare(middleware ...middleware.Middleware) {
	ep.MiddleWares = append(ep.MiddleWares, middleware...)
}

func (ep *Endpoint) UseObserve(observers ...middleware.Observer) {
	ep.Observers = append(ep.Observers, observers...)
}

func Find(name string) (*Endpoint, error) {
	for _, ep := range Endpoints {
		if ep.Name == name {
			return ep, nil
		}
	}
	return &Endpoint{}, errors.New("Location not found")
}

func (ep *Endpoint) RequestHTTP(value interface{}) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

func (ep *Endpoint) RequestStruct(value interface{}) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (u *WrappedURL) StringWithoutScheme() string {
	var buf bytes.Buffer
	if u.Opaque != "" {
		log.Printf("Writing opaque")
		buf.WriteString(u.Opaque)
	} else {
		if ui := u.User; ui != nil {
			buf.WriteString(ui.String())
			buf.WriteByte('@')
		}
		if h := u.Host; h != "" {
			log.Printf("We have a host.")
			buf.WriteString(h)
		}
		if u.Path != "" && u.Path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		buf.WriteString(escape(u.Path, 1))
	}

	if u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}

	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(escape(u.Fragment, 1))
	}

	return buf.String()
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
// When 'all' is true the full range of reserved characters are matched.
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last two as well. That leaves only ? to escape.
			return c == '?'

		case encodeUserPassword: // §3.2.2
			// The RFC allows ; : & = + $ , in userinfo, so we must escape only @ and /.
			// The parsing of userinfo treats : as special so we must escape that too.
			return c == '@' || c == '/' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}
