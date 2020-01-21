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
			glog.Fatal(err)
		}
		glog.V(1).Infof("%s\n", hdr.Name)
		//writer, found := contents.writers[hdr.Name]
		writer, found := contents.FindWriter(hdr.Name)
		if !found {
			glog.Warningf("unrecognized content in tar: %s", hdr.Name)
			continue
		}

		if _, err := io.Copy(writer, tr); err != nil {
			glog.Fatal(err)
		}
	}

	// validate loaded contents
	if err := contents.Validate(); err != nil {
		return err
	}

	return nil
}
