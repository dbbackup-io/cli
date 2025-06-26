package postgres

import (
	"context"
	"fmt"
	"io"
	"os"
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
		"-h", d.Host,
		"-p", fmt.Sprintf("%d", d.Port),
		"--format=custom",
		"--compress=6",
		"--no-password",
	}

	if d.Username != "" {
		args = append(args, "-U", d.Username)
	}

	if d.Database != "" {
		args = append(args, d.Database)
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)

	// Set password via environment variable if provided
	if d.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", d.Password))
	}

	// Create pipe for streaming output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to start pg_dump: %w", err)
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

// GetFileExtension returns the file extension for PostgreSQL dumps
func (d *Dumper) GetFileExtension() string {
	return ".dump"
}

// GetDatabaseType returns the database type
func (d *Dumper) GetDatabaseType() string {
	return "postgres"
}

// GetDatabaseName returns the database name
func (d *Dumper) GetDatabaseName() string {
	return d.Database
}