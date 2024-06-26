package cmd

import (
	"bytes"
	"errors"
	"sort"
	
	"github.com/olekukonko/tablewriter"
	sshtunnel "github.com/superwhys/venkit/ssh-tunnel/v2"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/vflags"
)

var (
	localAddr  string
	remoteAddr string
)

var (
	port     = vflags.Int("port", 0, "Port for serivce")
	env      = vflags.String("env", "", "Enviroment name for looking up connection profile")
	profiles = vflags.Struct("profiles", []*ConnectionProfile{}, "Connection profiles")
)

type ConnectionProfile struct {
	EnvName string
	Host    *sshtunnel.SshConfig
}

func (cp *ConnectionProfile) PopulateDefault(identityFile string) {
	if cp.Host.IdentityFile == "" {
		cp.Host.IdentityFile = identityFile
	}
	if cp.Host.User == "" {
		cp.Host.User = "root"
	}
}

func (cp *ConnectionProfile) Validate() error {
	if cp.Host.HostName == "" {
		return errors.New("SSH hostName requred")
	}
	return nil
}

func newTunnel() *sshtunnel.SshTunnel {
	var profile *ConnectionProfile
	var allProfiles []*ConnectionProfile
	lg.PanicError(profiles(&allProfiles))
	
	for _, p := range allProfiles {
		if p.EnvName == env() {
			profile = p
			break
		}
	}
	
	if profile == nil {
		lg.Fatal("no connection profile found. env= ", env())
	}
	
	return sshtunnel.NewTunnel(profile.Host)
}

// prettyMaps formats a map to table format string.
func prettyMaps(m map[string][]string) string {
	buffer := &bytes.Buffer{}
	table := tablewriter.NewWriter(buffer)
	table.SetColWidth(400)
	
	type Record struct {
		Uuid   string
		Typ    string
		Tunnel string
	}
	var rs []*Record
	for name, us := range m {
		r := &Record{
			Uuid:   name,
			Typ:    us[0],
			Tunnel: us[1],
		}
		
		rs = append(rs, r)
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Uuid < rs[j].Uuid
	})
	table.Append([]string{"UUID", "Tunnel-Type", "Tunnel"})
	for _, r := range rs {
		table.Append([]string{r.Uuid, r.Typ, r.Tunnel})
	}
	table.Render()
	return buffer.String()
}
