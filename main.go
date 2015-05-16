package main

import (
	"bytes"
	"errors"
	"flag"
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
	URL_PREFIX      = "http://met.hu/img/dewa"
	IMAGES_LOCATION = "images"
	GIF_LOCATION    = "gifs"
	GIF_FRAME_10MS  = 50
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

func EnsureDirectoryStructure() bool {
	err1 := ensureDirectory(IMAGES_LOCATION)
	err2 := ensureDirectory(GIF_LOCATION)
	if err1 != nil {
		log.Printf("Failed to create directory %v, error: %v", IMAGES_LOCATION, err1)
	}
	if err2 != nil {
		log.Printf("Failed to create directory %v, error: %v", GIF_LOCATION, err2)
	}

	return err1 == nil && err2 == nil
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
	result := make([]string, 24, 24)
	for hour := 0; hour < 24; hour++ {
		current := generateFileName(today, hour)
		result[hour] = current
	}
	return result
}

func isFileExist(fileName string) bool {
	stat, err := os.Stat(fileName)
	if os.IsNotExist(err) || err != nil || stat.Size() == 0 {
		return false
	}
	return true
}

func filterExistingFiles(l []string) []string {
	result := make([]string, 0, 0)
	for _, f := range l {
		if !isFileExist(imageName(f)) {
			result = append(result, f)
		}
	}
	return result
}

func imageName(f string) string {
	return fmt.Sprintf("./%v/%v", IMAGES_LOCATION, f)
}

func downloadFile(f string) error {
	//log.Println("Trying to download ", f)
	resp, err := http.Get(fmt.Sprintf("%v/%v", URL_PREFIX, f))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to download file, status: %v", resp.Status)
	}

	imageFile, err := os.Create(imageName(f))
	if err != nil {
		return err
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, resp.Body)
	return err
}

func tryDownloadFiles(l []string) {
	for _, f := range l {
		err := downloadFile(f)
		/*if err != nil {
			log.Printf("Failed to download %v: %v\n", f, err)
		} else {
			log.Printf("Downloaded %v\n", f)
		}*/
		if err == nil {
			log.Printf("Downloaded %v\n", f)
		}
	}
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

func saveGif(f string, img *gif.GIF) error {
	imgFile, err := os.Create(f)
	if err != nil {
		return err
	}

	err = gif.EncodeAll(imgFile, img)
	imgFile.Close()
	return err
}

func CreateGif(frameTime int, imagesLocation string, gifFileName string) error {
	imageFileInfos, err := ioutil.ReadDir(imagesLocation)
	if err != nil {
		return err
	}
	images := make([]*image.Paletted, 0, 0)
	delays := make([]int, 0, 0)
	// Load the images
	for _, f := range imageFileInfos {
		fileName := f.Name()
		im, err := loadImage(path.Join(imagesLocation, fileName))
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
	g := &gif.GIF{images, delays, 0 /*loop count*/}
	return saveGif(gifFileName, g)
}

func DownloadImages() {
	today := time.Now()
	l := generateFileNameList(today)
	l = filterExistingFiles(l)
	tryDownloadFiles(l)
}

func main() {
	log.Println(" --------- Started")
	doCheckStructure := flag.Bool("check", false, "Check directory structure")
	doDownloadImages := flag.Bool("download", false, "Download images")
	doCreateGif := flag.Bool("gif", false, "Generate gif file")

	flag.Parse()

	if *doCheckStructure {
		log.Println("Checking directory structure")
		ok := EnsureDirectoryStructure()
		if !ok {
			log.Fatal("Failed to create directory structure")
		}
	}

	if *doDownloadImages {
		log.Println("Downloading images")
		DownloadImages()
	}
	if *doCreateGif {
		log.Println("Creating gif")
		err := CreateGif(GIF_FRAME_10MS, IMAGES_LOCATION, path.Join(GIF_LOCATION, "anim.gif"))
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(" --------- Finished")
}
