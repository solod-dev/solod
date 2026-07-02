// Package net provides basic TCP, UDP, and Unix domain socket networking.
//
// It is a small subset of Go's net package. Supports TCP (networks "tcp",
// "tcp4", "tcp6"), UDP (networks "udp", "udp4", "udp6"), and Unix domain
// sockets ("unix" for streams and "unixgram" for datagrams).
//
// TCP is served by [ResolveTCPAddr], [DialTCP] and [ListenTCP] functions,
// and the associated [TCPConn] and [TCPListener] types.
//
// UDP is served by [ResolveUDPAddr], [DialUDP] (a connected socket, with
// [UDPConn.Read]/[UDPConn.Write]) and [ListenUDP] (an unconnected socket,
// with [UDPConn.ReadFrom]/[UDPConn.WriteTo]).
//
// Unix domain sockets are served by [ResolveUnixAddr], [DialUnix], [ListenUnix]
// (stream), and [ListenUnixgram] (datagram), sharing the [UnixConn] type;
// a [ListenUnix] or [ListenUnixgram] socket file is removed on Close.
//
// Accept, Read, and Write block by default. They can be bounded with
// a deadline: [TCPConn.SetDeadline], [TCPListener.SetDeadline],
// [UDPConn.SetDeadline], [UnixConn.SetDeadline] and [UnixListener.SetDeadline]
// make a pending call fail with ErrTimeout once the deadline passes. Without a
// deadline, a blocked call waits indefinitely.
//
// Unlike in Go, none of the types in this package are safe for concurrent use:
// a connection or listener must be used by one thread at a time. To serve
// connections concurrently, hand each accepted connection to its own thread
// (for example via conc.Pool).
package net

import (
	"solod.dev/so/bytealg"
)

// HostPort holds the host and port parts of an address.
type HostPort struct {
	Host string
	Port string
}

// SplitHostPort splits a network address of the form "host:port",
// "host%zone:port", "[host]:port" or "[host%zone]:port" into host or
// host%zone and port.
//
// The returned strings are views into hostport.
func SplitHostPort(hostport string) (HostPort, error) {
	i := bytealg.LastIndexByteString(hostport, ':')
	if i < 0 {
		return HostPort{}, ErrMissingPort
	}

	j, k := 0, 0
	var host, port string
	if hostport[0] == '[' {
		// Expect the first ']' just before the last ':'.
		end := bytealg.IndexByteString(hostport, ']')
		if end < 0 {
			return HostPort{}, ErrMissingBracket
		}
		switch end + 1 {
		case len(hostport):
			// There can't be a ':' behind the ']' now.
			return HostPort{}, ErrMissingPort
		case i:
			// The expected result.
		default:
			// Either ']' isn't followed by a colon, or it is
			// followed by a colon that is not the last one.
			if hostport[end+1] == ':' {
				return HostPort{}, ErrTooManyColons
			}
			return HostPort{}, ErrMissingPort
		}
		host = hostport[1:end]
		j, k = 1, end+1 // there can't be a '[' resp. ']' before these positions
	} else {
		host = hostport[:i]
		if bytealg.IndexByteString(host, ':') >= 0 {
			return HostPort{}, ErrTooManyColons
		}
	}

	if bytealg.IndexByteString(hostport[j:], '[') >= 0 {
		return HostPort{}, ErrUnexpectedBracket
	}
	if bytealg.IndexByteString(hostport[k:], ']') >= 0 {
		return HostPort{}, ErrUnexpectedBracket
	}

	port = hostport[i+1:]
	return HostPort{Host: host, Port: port}, nil
}

// JoinHostPort combines host and port into a network address of the
// form "host:port". If host contains a colon, as found in literal
// IPv6 addresses, then JoinHostPort returns "[host]:port".
//
// The result is built into buf and aliases it, so buf must have enough
// capacity (len(host)+len(port)+4 is always sufficient).
func JoinHostPort(buf []byte, host, port string) string {
	b := buf[:0]
	// We assume that host is a literal IPv6 address if host has colons.
	if bytealg.IndexByteString(host, ':') >= 0 {
		b = append(b, '[')
		b = append(b, host...)
		b = append(b, ']')
	} else {
		b = append(b, host...)
	}
	b = append(b, ':')
	b = append(b, port...)
	return string(b)
}
