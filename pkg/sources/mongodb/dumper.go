package mongodb

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type Dumper struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func (d *Dumper) CreateBackupStream(ctx context.Context) (io.ReadCloser, error) {
	// Build MongoDB connection URI
	uri := fmt.Sprintf("mongodb://%s:%d", d.Host, d.Port)
	if d.Username != "" && d.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", d.Username, d.Password, d.Host, d.Port)
	}

	args := []string{
		"--uri", uri,
		"--archive",
		"--gzip",
	}

	if d.Database != "" {
		args = append(args, "--db", d.Database)
	}

	cmd := exec.CommandContext(ctx, "mongodump", args...)

	// Create pipe for streaming output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to start mongodump: %w", err)
	}

	// Return a wrapper that will wait for the command to finish when closed
	return &cmdReader{
		reader: stdout,
		cmd:    cmd,
	}, nil
}

type cmdReader struct {
	reader io.ReadCloser
	cmd    *exec.Cmd
}

func (cr *cmdReader) Read(p []byte) (int, error) {
	return cr.reader.Read(p)
}

func (cr *cmdReader) Close() error {
	cr.reader.Close()
	return cr.cmd.Wait()
}

// GetFileExtension returns the file extension for MongoDB dumps
func (d *Dumper) GetFileExtension() string {
	return ".archive"
}

// GetDatabaseType returns the database type
func (d *Dumper) GetDatabaseType() string {
	return "mongodb"
}

// GetDatabaseName returns the database name
func (d *Dumper) GetDatabaseName() string {
	if d.Database == "" {
		return "all"
	}
	return d.Database
}
