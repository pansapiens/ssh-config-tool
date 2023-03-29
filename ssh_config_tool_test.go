package main

import (
	"reflect"
	"testing"
)

func TestParseSSHConfig(t *testing.T) {
	tests := []struct {
		name          string
		configData    string
		expectedHosts map[string][]string
	}{
		{
			name:          "Empty Config",
			configData:    "",
			expectedHosts: map[string][]string{},
		},
		{
			name: "Single Host",
			configData: `Host example
  HostName example.com
  User exampleuser
  Port 2222
  IdentityFile ~/.ssh/id_rsa`,
			expectedHosts: map[string][]string{
				"example": {
					"  HostName example.com",
					"  User exampleuser",
					"  Port 2222",
					"  IdentityFile ~/.ssh/id_rsa",
				},
			},
		},
		{
			name: "Multiple Hosts",
			configData: `Host example1
  HostName example1.com
  User user1

Host example2
  HostName example2.com
  User user2
  Port 2222`,
			expectedHosts: map[string][]string{
				"example1": {
					"  HostName example1.com",
					"  User user1",
				},
				"example2": {
					"  HostName example2.com",
					"  User user2",
					"  Port 2222",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hosts := parseSSHConfig(tt.configData)
			if !reflect.DeepEqual(hosts, tt.expectedHosts) {
				for host, lines := range hosts {
					if !reflect.DeepEqual(lines, tt.expectedHosts[host]) {
						t.Errorf("Mismatch in host %q:\nGot: %q\nWant: %q", host, lines, tt.expectedHosts[host])
					}
				}
			}
		})
	}
}

type SSHCommandToConfigEntryTestCase struct {
	Name          string
	SSHCommand    string
	ExpectedEntry string
}

type SSHCommandToConfigEntryTestFunc func(t *testing.T, tc SSHCommandToConfigEntryTestCase)

func TestSSHCommandToConfigEntryAlt(t *testing.T) {
	testCases := []SSHCommandToConfigEntryTestCase{
		{
			Name:          "Simple",
			SSHCommand:    "ssh user@host",
			ExpectedEntry: "Host host\n  HostName host\n  User user\n",
		},
		{
			Name:       "With Options",
			SSHCommand: `ssh -p 2222 -X -A -i ~/.ssh/somekey -L 5000:localhost:5000 -o "ProxyJump=bastionhost" -o "ForwardX11=yes" user@host`,
			ExpectedEntry: `Host host
  HostName host
  User user
  ForwardAgent yes
  ForwardX11 yes
  IdentityFile ~/.ssh/somekey
  LocalForward 5000 localhost:5000
  Port 2222
  ProxyJump bastionhost
`,
		},
	}

	runSSHCommandToConfigEntryTests(t, testCases, func(t *testing.T, tc SSHCommandToConfigEntryTestCase) {
		entry := sshCommandToConfigEntry(tc.SSHCommand)
		if entry != tc.ExpectedEntry {
			t.Errorf("expected:\n%v\n\ngot:\n%v\n", tc.ExpectedEntry, entry)
		}
	})
}

func runSSHCommandToConfigEntryTests(t *testing.T, testCases []SSHCommandToConfigEntryTestCase, testFunc SSHCommandToConfigEntryTestFunc) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testFunc(t, tc)
		})
	}
}
