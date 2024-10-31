module github.com/kjbreil/ziploc

go 1.23

toolchain go1.23.0

require (
	github.com/go-git/go-git/v5 v5.12.0
	github.com/go-ini/ini v1.67.0
	github.com/kjbreil/crcloc v0.0.0-20240619200433-6432ff988581
	github.com/kjbreil/crlf v0.0.0-20210116185654-a98352303dd9
	github.com/kjbreil/go-smb2 v1.1.0
	github.com/kjbreil/sil v0.0.0-20240619192433-4abc6ecd2d7a
	github.com/klauspost/compress v1.17.9
	github.com/schollz/progressbar/v3 v3.14.6
	github.com/spf13/cobra v1.8.1
	golang.org/x/text v0.18.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/cloudflare/circl v1.4.0 // indirect
	github.com/cyphar/filepath-securejoin v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/kjbreil/glsp v0.2.2 // indirect
	github.com/kjbreil/loc-macro v0.0.0-20240917185356-d96238df39ad // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/sourcegraph/jsonrpc2 v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

replace github.com/kjbreil/sil => ../sil

replace github.com/kjbreil/crcloc => ../crcloc

replace github.com/kjbreil/loc-macro => ../loc-macro

replace github.com/kjbreil/go-smb2 => ../go-smb2

replace github.com/kjbreil/glsp => ../glsp
