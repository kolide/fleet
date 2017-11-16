package okforward

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

const (
	defaultFastPort = 7651
)

type Writer struct {
	logger log.Logger
	buf    *bytes.Buffer
}

// New creates an io.Writer which forwards logs to oklog.
func New(logger log.Logger, peers []string) (*Writer, error) {
	urls, err := urls(peers)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := &Writer{
		logger: logger,
		buf:    buf,
	}

	go writer.loop(urls)

	return writer, nil
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *Writer) loop(urls []*url.URL) error {
	var prefix string
	// Shuffle the order.
	rand.Seed(time.Now().UnixNano())
	for i := range urls {
		j := rand.Intn(i + 1)
		urls[i], urls[j] = urls[j], urls[i]
	}

	backoff := time.Duration(0)

	// Enter the connect and forward loop. We do this forever.
	for ; ; urls = append(urls[1:], urls[0]) { // rotate thru URLs
		// We gonna try to connect to this first one.
		target := urls[0]

		host, port, err := net.SplitHostPort(target.Host)
		if err != nil {
			return errors.Wrapf(err, "unexpected error")
		}

		// Support e.g. "tcp+dnssrv://host:port"
		fields := strings.SplitN(target.Scheme, "+", 2)
		if len(fields) == 2 {
			proto, suffix := fields[0], fields[1]
			switch suffix {
			case "dns", "dnsip":
				ips, err := net.LookupIP(host)
				if err != nil {
					level.Info(w.logger).Log("LookupIP", host, "err", err)
					backoff = exponential(backoff)
					time.Sleep(backoff)
					continue
				}
				host = ips[rand.Intn(len(ips))].String()
				target.Scheme, target.Host = proto, net.JoinHostPort(host, port)

			case "dnssrv":
				_, records, err := net.LookupSRV("", proto, host)
				if err != nil {
					level.Info(w.logger).Log("LookupSRV", host, "err", err)
					backoff = exponential(backoff)
					time.Sleep(backoff)
					continue
				}
				host = records[rand.Intn(len(records))].Target
				target.Scheme, target.Host = proto, net.JoinHostPort(host, port) // TODO(pb): take port from SRV record?

			case "dnsaddr":
				names, err := net.LookupAddr(host)
				if err != nil {
					level.Info(w.logger).Log("LookupAddr", host, "err", err)
					backoff = exponential(backoff)
					time.Sleep(backoff)
					continue
				}
				host = names[rand.Intn(len(names))]
				target.Scheme, target.Host = proto, net.JoinHostPort(host, port)

			default:
				level.Info(w.logger).Log("unsupported_scheme_suffix", suffix, "using", proto)
				target.Scheme = proto // target.Host stays the same
			}
		}
		level.Debug(w.logger).Log("raw_target", urls[0].String(), "resolved_target", target.String())

		conn, err := net.Dial(target.Scheme, target.Host)
		if err != nil {
			level.Info(w.logger).Log("Dial", target.String(), "err", err)
			backoff = exponential(backoff)
			time.Sleep(backoff)
			continue
		}

		line, err := w.buf.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			level.Info(w.logger).Log("err", err)
			continue
		}
		record := string(bytes.TrimSpace(line))
		if n, err := fmt.Fprintf(conn, "%s%s\n", prefix, record); err != nil {
			level.Info(w.logger).Log("disconnected_from", target.String(), "due_to", err)
			break
		} else if n < len(record)+1 {
			level.Info(w.logger).Log("short_write_to", target.String(), "n", n, "less_than", len(record)+1)
			break // TODO(pb): we should do something more sophisticated here
		}

	}
	return nil
}

func urls(peers []string) ([]*url.URL, error) {
	// Parse URLs for forwarders.
	var urls []*url.URL
	for _, addr := range peers {
		schema, host, _, _, err := parseAddr(addr, defaultFastPort)
		if err != nil {
			return nil, errors.Wrap(err, "parsing ingest address")
		}
		u, err := url.Parse(fmt.Sprintf("%s://%s", schema, host))
		if err != nil {
			return nil, errors.Wrap(err, "parsing ingest URL")
		}
		if _, _, err := net.SplitHostPort(u.Host); err != nil {
			return nil, errors.Wrapf(err, "couldn't split host:port")
		}
		urls = append(urls, u)
	}
	return urls, nil
}

// "udp://host:1234", 80 => udp host:1234 host 1234
// "host:1234", 80       => tcp host:1234 host 1234
// "host", 80            => tcp host:80   host 80
func parseAddr(addr string, defaultPort int) (network, address, host string, port int, err error) {
	u, err := url.Parse(strings.ToLower(addr))
	if err != nil {
		return network, address, host, port, err
	}

	switch {
	case u.Scheme == "" && u.Opaque == "" && u.Host == "" && u.Path != "": // "host"
		u.Scheme, u.Opaque, u.Host, u.Path = "tcp", "", net.JoinHostPort(u.Path, strconv.Itoa(defaultPort)), ""
	case u.Scheme != "" && u.Opaque != "" && u.Host == "" && u.Path == "": // "host:1234"
		u.Scheme, u.Opaque, u.Host, u.Path = "tcp", "", net.JoinHostPort(u.Scheme, u.Opaque), ""
	case u.Scheme != "" && u.Opaque == "" && u.Host != "" && u.Path == "": // "tcp://host[:1234]"
		if _, _, err := net.SplitHostPort(u.Host); err != nil {
			u.Host = net.JoinHostPort(u.Host, strconv.Itoa(defaultPort))
		}
	default:
		return network, address, host, port, errors.Errorf("%s: unsupported address format", addr)
	}

	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		return network, address, host, port, err
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return network, address, host, port, err
	}

	return u.Scheme, u.Host, host, port, nil
}

func exponential(d time.Duration) time.Duration {
	const (
		min = 16 * time.Millisecond
		max = 1024 * time.Millisecond
	)
	d *= 2
	if d < min {
		d = min
	}
	if d > max {
		d = max
	}
	return d
}
