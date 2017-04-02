package kjudge

import "github.com/natsukagami/kjudge-api-go/lib/fs"

func folderClean(folder string) {
	fs.Remove(folder)
}
