package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const upPortFallbackRange = 20

// upPortPolicy distinguishes a user-requested listener from a default. Never
// change explicit ports silently; scripted users can enforce both defaults.
type upPortPolicy struct {
	strict          bool
	editorExplicit  bool
	previewExplicit bool
}

func (p upPortPolicy) fixed(service upService) bool {
	if p.strict {
		return true
	}
	switch service.name {
	case upServiceEditor:
		return p.editorExplicit
	case upServiceBuilder:
		return p.previewExplicit
	default:
		return true
	}
}

type upPortChange struct {
	service  string
	previous int
	selected int
}

// chooseUpPorts reserves no permanent listeners: the children own their
// sockets and validate them again when they start. Only loopback addresses are
// considered, and any race is reported by the child rather than stealing the
// port from an existing service.
func chooseUpPorts(ctx context.Context, services []upService, policy upPortPolicy) ([]upService, []upPortChange, error) {
	selected := append([]upService(nil), services...)
	claimed := make(map[int]string, len(services))
	changes := make([]upPortChange, 0, len(services))

	// Claim explicitly requested ports first. A fallback cannot take a
	// requested port, regardless of service order in the startup plan.
	for _, service := range selected {
		if !policy.fixed(service) {
			continue
		}
		port, err := upServicePort(service)
		if err != nil {
			return nil, nil, err
		}
		if owner, used := claimed[port]; used {
			return nil, nil, fmt.Errorf("%s port %d is already reserved by %s", service.name, port, owner)
		}
		if err := checkUpPorts(ctx, []upService{service}); err != nil {
			if isAddressInUseError(err) {
				return nil, nil, upOccupiedPortError(service, port)
			}
			return nil, nil, err
		}
		claimed[port] = service.name
	}

	for i, service := range selected {
		if policy.fixed(service) {
			continue
		}
		original, err := upServicePort(service)
		if err != nil {
			return nil, nil, err
		}
		available := false
		for candidate := original; candidate <= min(original+upPortFallbackRange, 65535); candidate++ {
			if err := ctx.Err(); err != nil {
				return nil, nil, err
			}
			if _, reserved := claimed[candidate]; reserved {
				continue
			}
			choice, err := withUpServicePort(service, candidate)
			if err != nil {
				return nil, nil, err
			}
			if err := checkUpPorts(ctx, []upService{choice}); err != nil {
				if isAddressInUseError(err) {
					continue
				}
				return nil, nil, err
			}
			selected[i] = choice
			claimed[candidate] = service.name
			if candidate != original {
				changes = append(changes, upPortChange{
					service: service.name, previous: original, selected: candidate,
				})
			}
			available = true
			break
		}
		if !available {
			return nil, nil, fmt.Errorf("%s cannot find a free loopback port in %d-%d; stop an existing service or rerun with --%s-port PORT",
				service.name, original, min(original+upPortFallbackRange, 65535), upPortFlag(service.name))
		}
	}
	return selected, changes, nil
}

func upServicePort(service upService) (int, error) {
	hostPort := strings.TrimSuffix(strings.TrimPrefix(service.url, "http://"), "/")
	host, portText, err := net.SplitHostPort(hostPort)
	if err != nil {
		return 0, fmt.Errorf("invalid %s listener %q: %w", service.name, service.url, err)
	}
	if host != "127.0.0.1" {
		return 0, fmt.Errorf("%s must listen on 127.0.0.1, got %q", service.name, host)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid %s port %q", service.name, portText)
	}
	return port, nil
}

func withUpServicePort(service upService, port int) (upService, error) {
	if _, err := upServicePort(service); err != nil {
		return upService{}, err
	}
	service.args = append([]string(nil), service.args...)
	service.url = "http://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(port)) + "/"
	flag, value := "", ""
	switch service.name {
	case upServiceEditor:
		flag, value = "--listen", net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	case upServiceBuilder:
		flag, value = "--port", strconv.Itoa(port)
	default:
		return upService{}, fmt.Errorf("unknown local service %q", service.name)
	}
	for i := 0; i+1 < len(service.args); i++ {
		if service.args[i] == flag {
			service.args[i+1] = value
			return service, nil
		}
	}
	return upService{}, fmt.Errorf("%s has no %s argument", service.name, flag)
}

func upPortFlag(service string) string {
	if service == upServiceEditor {
		return "editor"
	}
	return "preview"
}

func upOccupiedPortError(service upService, port int) error {
	flag := upPortFlag(service.name)
	return fmt.Errorf("%s port %d is already in use at %s; no process was stopped. "+
		"Inspect the existing listener or rerun with --%s-port PORT. "+
		"To use the next available local port automatically, omit --%s-port and --strict-ports",
		service.name, port, service.url, flag, flag)
}
