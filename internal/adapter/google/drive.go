package google

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
)

type GoogleDrive struct {
	Service *drive.Service
}

func (d *GoogleDrive) ListFiles(pageSize int) (map[string]string, error) {
	nameM := make(map[string]string)
	r, err := d.Service.Files.List().PageSize(int64(pageSize)).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		return nil, err
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			nameM[i.Id] = i.Name
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	fmt.Println("nameM:", nameM)
	return nameM, nil
}

func (d *GoogleDrive) ListMyDriveFolders() (map[string]string, error) {
	nameM := make(map[string]string)
	r, err := d.Service.Files.List().Q("mimeType='application/vnd.google-apps.folder'and 'root' in parents").
		Fields("files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		return nil, err
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			nameM[i.Id] = i.Name
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	fmt.Println("nameM 共:", len(nameM), "個資料夾")
	return nameM, nil
}

func (d *GoogleDrive) ListSharedFolders() (map[string]string, error) {
	nameM := make(map[string]string)
	r, err := d.Service.Files.List().Q("sharedWithMe and mimeType='application/vnd.google-apps.folder'").
		Fields("files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
		return nil, err
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			nameM[i.Id] = i.Name
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	fmt.Println("nameM 共:", len(nameM), "個資料夾")
	return nameM, nil
}
