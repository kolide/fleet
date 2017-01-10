package kolide

type FilePaths map[string][]string

type FileIntegrityMonitoringStore interface {
	NewFilePath(path *FilePath) (*FilePath, error)
	FilePaths() (FilePaths, error)
}

// FilePath maps a name to a group of files for the osquery file_paths
// section.
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
type FilePath struct {
	ID          uint
	SectionName string `db:"section_name"`
	Description string
	Paths       []string `db:"-"`
}
