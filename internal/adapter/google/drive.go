package google

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"os"
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

func (d *GoogleDrive) UploadFile(folderID string, fileName string, file *os.File) error {
	defer file.Close()
	// 指定目標資料夾的 ID
	var parents []string
	if folderID != "" {
		parents = []string{folderID}
	}

	// 上傳文件
	driveFile, err := d.Service.Files.Create(&drive.File{
		Name:    fileName,
		Parents: parents,
	}).Media(file).Do()
	if err != nil {
		log.Println("Upload Error:", err)
		return err
	}

	log.Printf("Got drive.File, err: %#v, %v", driveFile, err)
	return nil
}

func (d *GoogleDrive) ListFilesByID(folderID string) (map[string]string, error) {
	nameM := make(map[string]string)
	r, err := d.Service.Files.List().Q(fmt.Sprintf("'%s' in parents and mimeType != 'application/vnd.google-apps.folder'", folderID)).
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

	fmt.Println("nameM 共:", len(nameM), "個檔案")
	return nameM, nil
}

func (d *GoogleDrive) ListFolderByID(folderID string) (map[string]string, error) {
	nameM := make(map[string]string)
	r, err := d.Service.Files.List().Q(fmt.Sprintf("'%s' in parents and mimeType='application/vnd.google-apps.folder'", folderID)).
		Fields("files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve folders: %v", err)
		return nil, err
	}

	fmt.Println("Folders:")
	if len(r.Files) == 0 {
		fmt.Println("No folders found.")
	} else {
		for _, i := range r.Files {
			nameM[i.Id] = i.Name
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}

	fmt.Println("nameM 共:", len(nameM), "個資料夾")
	return nameM, nil
}

func (d *GoogleDrive) FindFolderPathByID(folderID string) (string, error) {
	return d.findFolderPathByIDRecursively(folderID, "")
}

func (d *GoogleDrive) findFolderPathByIDRecursively(folderID string, currentPath string) (string, error) {
	r, err := d.Service.Files.Get(folderID).Fields("id, name, parents,shared").Do()
	if err != nil {
		log.Printf("Unable to retrieve folder information: %v", err)
		return "", err
	}

	currentPath = r.Name + "/" + currentPath

	if len(r.Parents) > 0 {
		return d.findFolderPathByIDRecursively(r.Parents[0], currentPath)
	}

	if r.Shared {
		currentPath = "與我共用/" + currentPath
	}

	return currentPath, nil
}
