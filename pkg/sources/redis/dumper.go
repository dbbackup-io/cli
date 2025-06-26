package redis

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type Dumper struct {
	Host     string
	Port     int
	Password string
}

func (d *Dumper) CreateBackupStream(ctx context.Context) (io.ReadCloser, error) {
	args := []string{
		"-h", d.Host,
		"-p", fmt.Sprintf("%d", d.Port),
		"--rdb", "-",
	}

	if d.Password != "" {
		args = append(args, "-a", d.Password)
	}

	cmd := exec.CommandContext(ctx, "redis-cli", args...)

	// Create pipe for streaming output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to start redis-cli: %w", err)
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

// GetFileExtension returns the file extension for Redis dumps
func (d *Dumper) GetFileExtension() string {
	return ".rdb"
}

// GetDatabaseType returns the database type
func (d *Dumper) GetDatabaseType() string {
	return "redis"
}

// GetDatabaseName returns the database name (Redis doesn't have named databases like SQL)
func (d *Dumper) GetDatabaseName() string {
	return "default"
}