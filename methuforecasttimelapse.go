package methuforecasttimelapse

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	URL_PREFIX = "http://met.hu/img/dewa"
)

func ensureDirectory(dir string) error {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModeDir|0777)
	} else {
		if stat.IsDir() {
			return nil
		} else {
			return fmt.Errorf("File with name exists: %v", dir)
		}
	}
}

func generateFileName(t time.Time, hour int) string {
	return fmt.Sprintf(
		"dewa%04d%02d%02d_%02d00+Szeged.png",
		t.Year(),
		t.Month(),
		t.Day(),
		hour)
}

func generateFileNameList(today time.Time) []string {
	yesterday := today.Add(-1 * 24 * time.Hour)
	result := make([]string, 4, 4)
	result[0] = generateFileName(today, 0)
	result[1] = generateFileName(today, 12)
	result[2] = generateFileName(yesterday, 0)
	result[3] = generateFileName(yesterday, 12)
	return result
}

func isFileExist(fileName string) bool {
	stat, err := os.Stat(fileName)
	if os.IsNotExist(err) || err != nil || stat.Size() == 0 {
		return false
	}
	return true
}

func imageName(f, imageDir string) string {
	return fmt.Sprintf("./%v/%v", imageDir, f)
}

func filterExistingFiles(l []string, imageDir string) []string {
	result := make([]string, 0, 0)
	for _, f := range l {
		if !isFileExist(imageName(f, imageDir)) {
			result = append(result, f)
		}
	}
	return result
}

func downloadFile(f, imageDir string) error {
	resp, err := http.Get(fmt.Sprintf("%v/%v", URL_PREFIX, f))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to download file, status: %v", resp.Status)
	}

	imageFile, err := os.Create(imageName(f, imageDir))
	if err != nil {
		return err
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, resp.Body)
	return err
}

func tryDownloadFiles(l []string, imageDir string) (downloaded int) {
	for _, f := range l {
		err := downloadFile(f, imageDir)
		/*if err != nil {
			log.Printf("Failed to download %v: %v\n", f, err)
		} else {
			log.Printf("Downloaded %v\n", f)
		}*/
		if err == nil {
			log.Printf("Downloaded %v\n", f)
			downloaded += 1
		}
	}
	return
}

func loadImage(f string) (*image.Paletted, error) {
	imgFile, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return convertImage(img)
}

func convertImage(img image.Image) (*image.Paletted, error) {
	buf := bytes.Buffer{}
	// Encode the image as a gif to the buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	// Decode back as paletted image
	gifimg, err := gif.Decode(&buf)
	if err != nil {
		return nil, err
	}
	palettedImg, ok := gifimg.(*image.Paletted)
	if !ok {
		return nil, errors.New("Failed to convert image")
	}

	return palettedImg, nil
}

func saveGif(f string, img *gif.GIF) error {
	imgFile, err := os.Create(f)
	if err != nil {
		return err
	}

	err = gif.EncodeAll(imgFile, img)
	imgFile.Close()
	return err
}

func EnsureDirectoryStructure(imagesDir, gifDir string) bool {
	err1 := ensureDirectory(imagesDir)
	err2 := ensureDirectory(gifDir)
	if err1 != nil {
		log.Printf("Failed to create directory %v, error: %v", imagesDir, err1)
	}
	if err2 != nil {
		log.Printf("Failed to create directory %v, error: %v", gifDir, err2)
	}

	return err1 == nil && err2 == nil
}

func DownloadImages(imageDir string) (downloaded int) {
	today := time.Now()
	l := generateFileNameList(today)
	l = filterExistingFiles(l, imageDir)
	return tryDownloadFiles(l, imageDir)
}

func ListImageFiles(location string) ([]string, error) {
	imageFileInfos, err := ioutil.ReadDir(location)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(imageFileInfos), len(imageFileInfos))
	for i, f := range imageFileInfos {
		result[i] = path.Join(location, f.Name())
	}

	return result, nil
}

func CreateGif(frameTime int, imagesLocation string, gifFileName string) error {
	images := make([]*image.Paletted, 0, 0)
	delays := make([]int, 0, 0)
	// Load the images
	imageFiles, err := ListImageFiles(imagesLocation)
	if err != nil {
		return err
	}
	for _, fileName := range imageFiles {
		im, err := loadImage(fileName)
		if err == nil {
			delays = append(delays, frameTime)
			images = append(images, im)
		} else {
			log.Printf("Failed to load image '%v', %v\n", fileName, err)
		}
	}
	if len(images) == 0 {
		return fmt.Errorf("Failed to load any images from '%v'", imagesLocation)
	}
	// Create and save the GIF
	g := &gif.GIF{Image: images, Delay: delays, LoopCount: 0}
	return saveGif(gifFileName, g)
}
