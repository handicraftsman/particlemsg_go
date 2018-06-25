package particlemsg

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// LoadPlugins - loads plugins listed in ClientConfig
func LoadPlugins(host, port string, clients *ClientConfig) {
	for _, client := range *clients {
		if !client.DoNotLoad {
			go func() {
				starter := func() {
					p, err := filepath.Abs(client.Path)
					if err != nil {
						panic(err)
					}

					d := filepath.Dir(p)

					cmd := exec.Command(client.Path)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Dir = d
					cmd.Env = append(os.Environ(),
						"PMSG_HOST="+host,
						"PMSG_PORT="+port,
						"PMSG_NAME="+client.Name,
						"PMSG_KEY="+client.Key,
					)
					if err = cmd.Run(); err != nil {
						panic(err)
					}
					time.Sleep(time.Second * 10)
				}
				for {
					starter()
				}
			}()
		}
	}
}
