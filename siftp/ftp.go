package siftp

import (
	"bytes"
	"crypto/tls"
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

type Client struct {
	Addr         string
	userID       string
	userPW       string
	connTimeout  time.Duration
	writeTimeout time.Duration
	readTimeout  time.Duration
	tlsConfig    *tls.Config
	epsvEnabled  bool
}

func NewClient(addr string, userID, userPW string) *Client {
	return &Client{
		Addr:         addr,
		userID:       userID,
		userPW:       userPW,
		connTimeout:  6 * time.Second,
		writeTimeout: 6 * time.Second,
		readTimeout:  6 * time.Second,
		tlsConfig:    nil,
		epsvEnabled:  false,
	}
}

func (ftpClient *Client) ReadFile(fileName string) ([]byte, error) {
	c, err := ftp.Dial(ftpClient.Addr,
		ftp.DialWithTimeout(ftpClient.connTimeout),
		ftp.DialWithDisabledEPSV(ftpClient.epsvEnabled))

	if err != nil {
		return nil, err
	}

	if err := c.Login(ftpClient.userID, ftpClient.userPW); err != nil {
		return nil, err
	}

	r, err := c.Retr(fileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	r.SetDeadline(time.Now().Add(ftpClient.readTimeout))
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err := c.Quit(); err != nil {
		return nil, err
	}

	return buf, nil
}

func (ftpClient *Client) WriteFile(fileName string, data []byte) error {
	c, err := ftp.Dial(ftpClient.Addr,
		ftp.DialWithTimeout(ftpClient.connTimeout),
		ftp.DialWithDisabledEPSV(ftpClient.epsvEnabled))

	if err != nil {
		return err
	}

	if err := c.Login(ftpClient.userID, ftpClient.userPW); err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	if err := c.Stor(fileName, buf); err != nil {
		return err
	}

	if err := c.Quit(); err != nil {
		return err
	}

	return nil
}
