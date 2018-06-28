package particlemsg

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// LoadPlugin - loads given plugin
func LoadPlugin(host, port, crt, key string, client *ClientInfo) {
	starter := func() {
		p, err := filepath.Abs(client.Path)
		if err != nil {
			panic(err)
		}

		d := filepath.Dir(p)

		var unsafeSSL string
		if client.UnsafeSSL {
			unsafeSSL = "true"
		} else {
			unsafeSSL = "false"
		}

		crt, _ = filepath.Abs(crt)
		key, _ = filepath.Abs(key)

		cmd := exec.Command(p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = d
		cmd.Env = append(os.Environ(),
			"PMSG_HOST="+host,
			"PMSG_PORT="+port,
			"PMSG_NAME="+client.Name,
			"PMSG_KEY="+client.Key,
			"PMSG_UNSAFE_SSL="+unsafeSSL,
			"PMSG_SSL_CERT="+crt,
			"PMSG_SSL_KEY="+key,
		)
		cmd.Env = append(cmd.Env, client.Env...)
		if err = cmd.Run(); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 10)
	}
	for {
		starter()
	}
}

// LoadPlugins - loads plugins listed in ClientConfig
func LoadPlugins(host, port, crt, key string, clients *ClientConfig, doNotLoadCore bool) {
	for _, client := range *clients {
		if !client.DoNotLoad {
			if doNotLoadCore && client.Name == "core" {
				continue
			}
			go LoadPlugin(host, port, crt, key, &client)
		}
	}
}

// FindClientInfo - finds client info by plugin name
func FindClientInfo(clients *ClientConfig, name string) (*ClientInfo, error) {
	for _, client := range *clients {
		if client.Name == name {
			return &client, nil
		}
	}
	return nil, errors.New("cannot find " + name)
}
