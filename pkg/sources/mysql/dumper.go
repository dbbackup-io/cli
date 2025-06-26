package mysql

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
	args := []string{
		fmt.Sprintf("--host=%s", d.Host),
		fmt.Sprintf("--port=%d", d.Port),
		"--single-transaction",
		"--routines",
		"--triggers",
	}

	if d.Username != "" {
		args = append(args, fmt.Sprintf("--user=%s", d.Username))
	}

	if d.Password != "" {
		args = append(args, fmt.Sprintf("--password=%s", d.Password))
	}

	if d.Database != "" {
		args = append(args, d.Database)
	}

	cmd := exec.CommandContext(ctx, "mysqldump", args...)

	// Create pipe for streaming output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to start mysqldump: %w", err)
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

// GetFileExtension returns the file extension for MySQL dumps
func (d *Dumper) GetFileExtension() string {
	return ".sql"
}

// GetDatabaseType returns the database type
func (d *Dumper) GetDatabaseType() string {
	return "mysql"
}

// GetDatabaseName returns the database name
func (d *Dumper) GetDatabaseName() string {
	return d.Database
}