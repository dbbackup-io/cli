package postgres

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// Create pipes for both stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdout.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdout.Close()
		stderr.Close()
		return nil, fmt.Errorf("failed to start pg_dump: %w", err)
	}

	// Return a wrapper that will wait for the command to finish when closed
	return &cmdReader{
		reader:    stdout,
		stderr:    stderr,
		cmd:       cmd,
		validated: false,
	}, nil
}

type cmdReader struct {
	reader    io.ReadCloser
	stderr    io.ReadCloser
	cmd       *exec.Cmd
	validated bool
	firstRead bool
}

func (cr *cmdReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)

	// On first read, validate that we got actual data
	if !cr.firstRead {
		cr.firstRead = true

		// If we got no data and EOF immediately, pg_dump likely failed
		if n == 0 && err == io.EOF {
			// Wait briefly for process to exit and check status
			time.Sleep(100 * time.Millisecond)

			// Check stderr for error messages
			stderrBytes, _ := io.ReadAll(cr.stderr)
			if len(stderrBytes) > 0 {
				errMsg := string(stderrBytes)
				// If stderr contains error messages, pg_dump failed
				if strings.Contains(errMsg, "error:") ||
					strings.Contains(errMsg, "FATAL:") ||
					strings.Contains(errMsg, "authentication failed") {
					return 0, fmt.Errorf("pg_dump failed: %s", errMsg)
				}
			}

			// Also check process state
			if cr.cmd.ProcessState != nil && !cr.cmd.ProcessState.Success() {
				if len(stderrBytes) > 0 {
					return 0, fmt.Errorf("pg_dump failed: %s", string(stderrBytes))
				}
				return 0, fmt.Errorf("pg_dump failed with exit code: %d", cr.cmd.ProcessState.ExitCode())
			}
		} else if n > 0 {
			cr.validated = true
		}
	}

	return n, err
}

func (cr *cmdReader) Close() error {
	cr.reader.Close()

	// Read remaining stderr to capture any error messages
	remainingStderr, _ := io.ReadAll(cr.stderr)
	cr.stderr.Close()

	if err := cr.cmd.Wait(); err != nil {
		if len(remainingStderr) > 0 {
			// Check if we should log debug output
			stderrStr := string(remainingStderr)
			if strings.Contains(os.Getenv("LOG_LEVEL"), "debug") ||
				strings.Contains(os.Getenv("LOG_LEVEL"), "DEBUG") {
				log.Printf("DEBUG: pg_dump stderr: %s", stderrStr)
			}
			return fmt.Errorf("pg_dump failed: %w\nOutput: %s", err, stderrStr)
		}
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	// Also log any warnings/messages even on success if debug is enabled
	if len(remainingStderr) > 0 {
		if strings.Contains(os.Getenv("LOG_LEVEL"), "debug") ||
			strings.Contains(os.Getenv("LOG_LEVEL"), "DEBUG") {
			log.Printf("DEBUG: pg_dump stderr: %s", string(remainingStderr))
		}
	}

	return nil
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
