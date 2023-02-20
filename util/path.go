package util

//import "errors"
import "os"

func GetCurrentDir() (string, error) {
	return os.Getwd()
}

func GetMonDataDir() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := d + "/mondata"
	return dir, nil
}

func GetDownloadDir() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := d + "/download"
	return dir, nil
}

func GetResDir() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := d + "/res"
	return dir, nil
}

func IsFileExist(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}
