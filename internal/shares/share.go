package shares

import (
	"fmt"
	"github.com/kjbreil/crlf"
	"github.com/kjbreil/go-smb2"
	"io"
	iofs "io/fs"
	"net"
	"os"
	"regexp"
	"strings"
)

type Share struct {
	path     string
	user     string
	password string
	server   string
	share    string
	conn     net.Conn
	session  *smb2.Session
	fs       *smb2.Share
}

func New(path, user, password string) (*Share, error) {

	parts := regexp.MustCompile(`(?m)\\\\(?<server>[^\\]+)\\(?<share>[^\\]+)$`).FindStringSubmatch(path)

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid share path: %s", path)
	}

	return &Share{
		path:     path,
		server:   parts[1],
		share:    parts[2],
		user:     user,
		password: password,
	}, nil
}

func (s *Share) Connect() error {
	var err error
	s.conn, err = net.Dial("tcp", fmt.Sprintf("%s:445", s.server))
	if err != nil {
		return err
	}
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.user,
			Password: s.password,
		},
	}

	s.session, err = d.Dial(s.conn)
	if err != nil {
		return err
	}

	s.fs, err = s.session.Mount(s.share)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *Share) Close() {
	s.fs.Umount()
	s.session.Logoff()
	s.conn.Close()
}

func (s *Share) WalkDir(path string, onFile func(path string, d iofs.DirEntry) error) error {

	if path == "" {
		path = "."
	}

	if strings.HasPrefix(path, s.path) {
		path = path[len(s.path):]
	}

	err := iofs.WalkDir(s.fs.DirFS(path), path, func(walkPath string, d iofs.DirEntry, err error) error {
		// stat, err := s.fs.Stat(path)
		if err != nil {
			return err
		}

		err = onFile(walkPath, d)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Share) ReadFile(path string) ([]byte, error) {
	if strings.HasPrefix(path, s.path) {
		path = path[len(s.path):]
	}
	if strings.HasPrefix(path, "\\") {
		path = path[1:]
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	file, err := s.fs.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Share) WriteFile(path string, data []byte) error {
	if strings.HasPrefix(path, s.path) {
		path = path[len(s.path):]
	}
	if strings.HasPrefix(path, "\\") {
		path = path[1:]
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	file, err := s.fs.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	nwc := writerCloser(file)

	_, err = nwc.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Share) SetArchive(fn string, archive bool) error {

	return s.fs.SetArchive(fn, archive)
	// return s.fs.Chmod(fn, 0444)
}

func (s *Share) Stat(path string) (iofs.DirEntry, error) {
	info, err := s.fs.Stat(path)
	if err != nil {
		return nil, err
	}
	return iofs.FileInfoToDirEntry(info), nil
}

type wc struct {
	io.Writer
	io.Closer
}

func writerCloser(f io.WriteCloser) io.WriteCloser {

	return wc{
		Writer: crlf.NewWriter(f),
		Closer: f,
	}
}
