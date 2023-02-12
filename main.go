package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	GithubOwner      string = "goeno"
	GithubRepository string = "aum"
	DS                      = string(os.PathSeparator)
)

var AppdataFolderPath string
var AppFolderPath string
var DownloadFolderPath string
var BinFolderPath string

func init() {
	var err error
	AppdataFolderPath, err = os.UserConfigDir()
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
	AppFolderPath = fmt.Sprintf("%s%saum", AppdataFolderPath, DS)
	mkdirIfNoExists(AppFolderPath)
	DownloadFolderPath = fmt.Sprintf("%s%sdownloads", AppFolderPath, DS)
	mkdirIfNoExists(DownloadFolderPath)
	BinFolderPath = fmt.Sprintf("%s%sbin", AppFolderPath, DS)
	mkdirIfNoExists(BinFolderPath)
}

func main() {
	if !amAdmin() {
		ShowError(errors.New("AmongUsMods updater require a terminal running as administrator. Please reopen your terminal or this updater with administrator permissions"))
		Stop()
		return
	}
	executable, err := os.Executable()
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
	updaterPath := fmt.Sprintf("%s%saum-update.exe", BinFolderPath, DS)
	if !fileExists(updaterPath) {
		_, err := copyFile(executable, updaterPath)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
	} else {
		currentUpdaterMd5, err := getFileMd5(executable)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
		targetUpdaterMd5, err := getFileMd5(executable)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
		if currentUpdaterMd5 != targetUpdaterMd5 {
			_, err := copyFile(executable, updaterPath)
			if err != nil {
				ShowError(err)
				Stop()
				return
			}
		}
	}
	release, err := GetLatestRelease(GithubOwner, GithubRepository)
	if err != nil {
		ShowError(err)
		Stop()
		return
	}
	if release == nil {
		ShowError(errors.New("unable to find valid release"))
		Stop()
		return
	}
	asset := release.GetExeAsset()
	if asset == nil {
		ShowError(errors.New("unable to find valid asset"))
		Stop()
		return
	}
	DownloadFile(asset.DownloadUrl, DownloadFolderPath)
	filePath := fmt.Sprintf("%s%s%s", DownloadFolderPath, DS, asset.Name)
	targetFilePath := fmt.Sprintf("%s%saum.exe", BinFolderPath, DS)
	pathEnv := os.Getenv("PATH")
	if strings.Contains(pathEnv, BinFolderPath) {
		fmt.Println("oui")
	} else {
		pathEnv += fmt.Sprintf(";%s", BinFolderPath)
		err = os.Setenv("PATH", pathEnv)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
	}
	if !fileExists(fmt.Sprintf("%s%saum.exe", BinFolderPath, DS)) {
		_, err = copyFile(filePath, targetFilePath)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
		ShowSuccess("AmongUsMods successfully installed")
	} else {
		target, err := getFileMd5(targetFilePath)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
		current, err := getFileMd5(filePath)
		if err != nil {
			ShowError(err)
			Stop()
			return
		}
		if target != current {
			_, err = copyFile(filePath, targetFilePath)
			if err != nil {
				ShowError(err)
				Stop()
				return
			}
			ShowSuccess("AmongUsMods successfully updated")
		} else {
			ShowInfo("AmongUsMods is already up to date")
		}
	}
	err = os.Remove(fmt.Sprintf("%s/%s", DownloadFolderPath, asset.Name))
	if err != nil {
		ShowError(err)
	}
	Stop()
}
