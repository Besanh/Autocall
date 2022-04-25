package sftp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type (
	Config struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	SFTPClient struct {
		sync.Mutex
		sshConn    *ssh.Client
		sftpClient *sftp.Client
		shutdown   chan bool
		closed     bool
		reconnects uint64
	}
	ISFTPManager interface {
		NewClient() (*SFTPClient, error)
		GetConnection() (*SFTPClient, error)
		SetLogger(logger *log.Logger)
		Close() error
	}
	SFTPManager struct {
		conns      []*SFTPClient
		log        *log.Logger
		connString string
		sshConfig  *ssh.ClientConfig
	}
)

var SFTPConfig Config
var BasicSFTPManager *SFTPManager

func NewSFTPManager(config Config) (*SFTPManager, error) {
	switch {
	case strings.TrimSpace(config.Host) == "",
		strings.TrimSpace(config.User) == "",
		strings.TrimSpace(config.Password) == "",
		0 >= config.Port || config.Port > 65535:
		return nil, errors.New("invalid parameters")
	}
	connString := fmt.Sprintf("%s:%d", config.Host, config.Port)
	SFTPConfig = config
	sshConfig := &ssh.ClientConfig{
		User:    config.User,
		Timeout: 30 * time.Second,
		// Auth: []ssh.AuthMethod{
		// 	ssh.KeyboardInteractive(SshInteractive),
		// },
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	manager := &SFTPManager{
		conns:      make([]*SFTPClient, 0),
		connString: connString,
		sshConfig:  sshConfig,
		log:        log.New(),
	}
	return manager, nil
}
func (sm *SFTPManager) NewClient() (client *SFTPClient, err error) {
	conn, err := ssh.Dial("tcp", sm.connString, sm.sshConfig)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	client = &SFTPClient{
		sshConn:    conn,
		sftpClient: sftpClient,
		shutdown:   make(chan bool, 1),
	}
	go sm.handleReconnects(client)
	sm.conns = append(sm.conns, client)
	return client, nil
}

func (sm *SFTPManager) handleReconnects(c *SFTPClient) {
	closed := make(chan error, 1)
	go func() {
		closed <- c.sshConn.Wait()
	}()

	select {
	case <-c.shutdown:
		c.sshConn.Close()
		break
	case res := <-closed:
		sm.log.Printf("Connection closed, reconnecting: %s", res)
		conn, err := ssh.Dial("tcp", sm.connString, sm.sshConfig)
		if err != nil {
			log.Error("Failed to reconnect:" + err.Error())
			break
		}

		SFTPClient, err := sftp.NewClient(conn)
		if err != nil {
			log.Error("Failed to reconnect:" + err.Error())
			break
		}

		atomic.AddUint64(&c.reconnects, 1)
		c.Lock()
		c.sftpClient = SFTPClient
		c.sshConn = conn
		c.Unlock()
		// Cool we have a new connection, keep going
		sm.handleReconnects(c)
	}
}

// Close closes the underlying connections
func (s *SFTPClient) Close() error {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return errors.New("connection was already closed")
	}

	s.shutdown <- true
	s.closed = true
	s.sshConn.Close()
	return s.sshConn.Wait()
}

// GetClient returns the underlying *sftp.Client
func (s *SFTPClient) GetClient() *sftp.Client {
	s.Lock()
	defer s.Unlock()
	return s.sftpClient
}

// Upload file to sftp server
func (sc *SFTPClient) Put(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		err := sc.sftpClient.Mkdir(path)
		_ = err
	}

	dstFile, err := sc.sftpClient.Create(remoteFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}
func SshInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	// The second parameter is unused
	for n := range questions {
		answers[n] = SFTPConfig.Password
	}

	return answers, nil
}

// GetConnection returns one of the existing connections the manager knows about. If there
// is no connections, we create a new one instead.
func (sm *SFTPManager) GetConnection() (*SFTPClient, error) {
	if len(sm.conns) > 0 {
		return sm.conns[0], nil
	}
	return sm.NewClient()
}

// SetLogger allows you to override the logger
func (sm *SFTPManager) SetLogger(logger *log.Logger) {
	sm.log = logger
}

// Close closes all connections managed by this manager
func (sm *SFTPManager) Close() error {
	for _, c := range sm.conns {
		c.Close()
	}
	return nil
}
