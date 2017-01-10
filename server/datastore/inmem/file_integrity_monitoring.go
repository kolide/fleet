package inmem

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewFilePath(fp *kolide.FilePath) (*kolide.FilePath, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	fp.ID = d.nextID(fp)
	d.filePaths[fp.ID] = fp
	return fp, nil
}

func (d *Datastore) FilePaths() (kolide.FilePaths, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	result := make(kolide.FilePaths)
	for _, filePath := range d.filePaths {
		result[filePath.SectionName] = append(result[filePath.SectionName], filePath.Paths...)
	}
	return result, nil
}
