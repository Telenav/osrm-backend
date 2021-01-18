package osrmfiles

import (
	"archive/tar"
	"io"
	"os"

	"github.com/golang/glog"
)

// Load `.osrm[.xxx]` tar file contents into its Contents structure.
func Load(contents ContentsOperator) error {
	f, err := os.Open(contents.FilePath())
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(2).Infof("open %s succeed.\n", contents.FilePath())

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			glog.Errorf("parsing tar file %s failed, err: %v", contents.FilePath(), err)
			break
		}
		glog.V(1).Infof("%s\n", hdr.Name)
		//writer, found := contents.writers[hdr.Name]
		writer, found := contents.FindWriter(hdr.Name)
		if !found {
			glog.Warningf("parsing tar file %s but content unrecognized: %s", contents.FilePath(), hdr.Name)
			continue
		}

		// The `io.Write()` will be called many times in `io.Copy()` due to fixed buffer size.
		// Defaultly the copy buffer size is 32*1024 bytes, see line 391 `copyBuffer()` in https://golang.org/src/io/io.go.
		// Unfortunetly, sometimes we expect to parse the packed data which can not be divided by 32*1024, e.g., 32*1024 / 12 != 0.
		// So we simply to use the Greatest common divisor to solve the issue.
		buf := make([]byte, 12*32*1024)
		if _, err := io.CopyBuffer(writer, tr, buf); err != nil {
			glog.Errorf("parsing content %s from tar file %s failed, err: %v", hdr.Name, contents.FilePath(), err)
			continue
		}
	}

	// post process once all contents loaded
	if err := contents.PostProcess(); err != nil {
		return err
	}

	// validate loaded contents
	if err := contents.Validate(); err != nil {
		return err
	}

	return nil
}
