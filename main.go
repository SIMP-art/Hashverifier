package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("\n")
		fmt.Printf("Usage: %s <filename>", os.Args[0])
		return
	}
	if os.Args[1] != "compare" {
		if len(os.Args) < 2 {
			fmt.Printf("Usage: go run %s <filename>", os.Args[0])
			fmt.Printf("\n")
			fmt.Printf("Usage: %s <filename>", os.Args[0])
			return
		}
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		md5Hasher := md5.New()
		if _, err := io.Copy(md5Hasher, file); err != nil {
			log.Fatal(err)
		}
		md5Hash := md5Hasher.Sum(nil)
		md5HashString := hex.EncodeToString(md5Hash)
		sha1Hasher := sha1.New()
		if _, err := file.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(sha1Hasher, file); err != nil {
			log.Fatal(err)
		}
		sha1Hash := sha1Hasher.Sum(nil)
		sha1HashString := hex.EncodeToString(sha1Hash)
		sha256Hasher := sha256.New()
		if _, err := file.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(sha256Hasher, file); err != nil {
			log.Fatal(err)
		}
		sha256Hash := sha256Hasher.Sum(nil)
		sha256HashString := hex.EncodeToString(sha256Hash)

		fileType, err := GetFileType(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		if fileType == "x-msdownload" {
			fileType = "window-excutable"
		}
		fmt.Printf("\n          *=================== %s | %s ===================*\n", os.Args[1], fileType)
		res("Sha1sum", "Sha256sum", "MD5sum", sha1HashString, sha256HashString, md5HashString)
	}
	if len(os.Args) != 4 {
		fmt.Printf("use compare mode: %s compare <file1> <file2>", os.Args[0])
		return
	}
	filematch, err, hash_compare := Compare(os.Args[2], os.Args[3])
	if err != nil {
		//fmt.Printf(err)
	}

	if filematch {
		res(os.Args[2], os.Args[3], "same?", hash_compare[0], hash_compare[1], "Yes")
	} else {
		res(os.Args[2], os.Args[3], "same?", hash_compare[0], hash_compare[1], "No")
	}
}
func GetFileType(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	fileType := mime.TypeByExtension(filepath.Ext(filename))
	if fileType == "" {
		fileType = http.DetectContentType(buffer)
	}
	parts := strings.Split(fileType, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("unknown file type")
	}
	extension := parts[1]

	return extension, nil
}

func res(h1, h2, h3, sha1, sha256, md5 string) {
	data := [][]string{
		{h1, h2, h3},
		{sha1, sha256, md5},
	}
	columnWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}
	border := "╔" + strings.Repeat("═", columnWidths[0]+2) + "╦" + strings.Repeat("═", columnWidths[1]+2) + "╦" + strings.Repeat("═", columnWidths[2]+2) + "╗"
	fmt.Println(border)
	for i, row := range data {
		fmt.Printf("║ %-*s ║ %-*s ║ %-*s ║\n", columnWidths[0], row[0], columnWidths[1], row[1], columnWidths[2], row[2])
		if i == 0 {
			fmt.Printf("╠%s╬%s╬%s╣\n", strings.Repeat("═", columnWidths[0]+2), strings.Repeat("═", columnWidths[1]+2), strings.Repeat("═", columnWidths[2]+2))
		}
		if i == len(data)-1 {
			border = "╚" + strings.Repeat("═", columnWidths[0]+2) + "╩" + strings.Repeat("═", columnWidths[1]+2) + "╩" + strings.Repeat("═", columnWidths[2]+2) + "╝"
			fmt.Println(border)
		}
	}
}

func Compare(file1, file2 string) (bool, error, []string) {
	//fmt.Println(file1, file2)
	data1, err := ioutil.ReadFile(file1)
	if err != nil {
		return false, err, []string{}
	}
	data2, err := ioutil.ReadFile(file2)
	if err != nil {
		return false, err, []string{}
	}
	hash1 := sha256.Sum256(data1)
	hash1Str := hex.EncodeToString(hash1[:])
	hash2 := sha256.Sum256(data2)
	hash2Str := hex.EncodeToString(hash2[:])
	return hash1Str == hash2Str, nil, []string{hash1Str, hash2Str}
}
