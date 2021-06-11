package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/sura8257/skbn/pkg/skbn"
)

func main() {
	cmd := NewRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		log.Fatal("Failed to execute command")
	}
}

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skbn",
		Short: "",
		Long:  ``,
	}

	out := cmd.OutOrStdout()

	cmd.AddCommand(NewCpCmd(out))
	cmd.AddCommand(NewVersionCmd(out))

	return cmd
}

type cpCmd struct {
	src        string
	dst        string
	parallel   int
	bufferSize int64

	out io.Writer
}

// NewCpCmd represents the copy command
func NewCpCmd(out io.Writer) *cobra.Command {
	c := &cpCmd{out: out}

	cmd := &cobra.Command{
		Use:   "cp",
		Short: "Copy a file to and from cloud storage",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := skbn.Copy(c.src, c.dst, c.parallel, c.bufferSize); err != nil {
				log.Fatal(err)
			}
		},
	}
	f := cmd.Flags()

	f.StringVar(&c.src, "src", "", "path to copy from. Example: path/to/copyfrom")
	f.StringVar(&c.dst, "dst", "", "path to copy to. Example: s3://<bucketName>/path/to/copyto")
	f.IntVarP(&c.parallel, "parallel", "p", 0, "number of parallel per call to upload when sending parts. If this is set to zero, the DefaultUploadConcurrency value will be used")
	f.Int64VarP(&c.bufferSize, "buffer-size", "b", 0, "The buffer size to use when buffering data into chunks and sending them as parts to S3. If this value is set to zero, the DefaultUploadPartSize value will be used.")

	cmd.MarkFlagRequired("src")
	cmd.MarkFlagRequired("dst")

	return cmd
}

var (
	// GitTag stands for a git tag
	GitTag string
	// GitCommit stands for a git commit hash
	GitCommit string
)

// NewVersionCmd prints version information
func NewVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version %s (git-%s)\n", GitTag, GitCommit)
		},
	}

	return cmd
}
