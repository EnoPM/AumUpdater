package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path"
)

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
}

func closeFile(out *os.File) {
	err := out.Close()
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
}

func ShowError(err error) {
	fmt.Printf("%s %s", RedBold("[ERROR]"), Red(err.Error()))
}

func ShowSuccess(message string) {
	fmt.Printf("%s %s", GreenBold("[SUCCESS]"), Green(message))
}

func ShowInfo(message string) {
	fmt.Printf("%s %s", BlueBold("[INFO]"), Blue(message))
}

func Stop() {
	fmt.Print("\nPress 'Enter' to close updater...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Red(text string) string {
	return color.RedString(text)
}

func Green(text string) string {
	return color.GreenString(text)
}

func Blue(text string) string {
	return color.BlueString(text)
}

func RedBold(text string) string {
	c := color.New(color.Bold, color.FgRed)
	return c.Sprint(text)
}

func GreenBold(text string) string {
	c := color.New(color.Bold, color.FgGreen)
	return c.Sprint(text)
}

func BlueBold(text string) string {
	c := color.New(color.Bold, color.FgBlue)
	return c.Sprint(text)
}

func DownloadFile(url string, dest string) {
	fileName := path.Base(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
	defer closeBody(resp.Body)
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", dest, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	defer closeFile(f)
	if err != nil {
		ShowError(err)
		Stop()
		return
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fileName,
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
}

func mkdirIfNoExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			ShowError(err)
		}
	}
}

func amAdmin() bool {
	file, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	closeFile(file)
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer closeFile(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer closeFile(destination)
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func getFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer closeFile(file)
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
