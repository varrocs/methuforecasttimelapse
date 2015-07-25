package methuforecasttimelapse

import (
	"flag"
	_ "image/png"
	"io"
	"log"
	"os"
	"path"
	"time"
)

const (
	LOG_FILE_NAME   = "log.txt"
	CHECKING_PERIOD = 1 * time.Hour
	//CHECKING_PERIOD = 10 * time.Second
)

func setupLogger() *os.File {
	logFile, err := os.OpenFile(LOG_FILE_NAME, os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		return logFile
	} else {
		log.Println(err)
		return nil
	}
}

func main() {
	//logFile := setupLogger()
	//defer logFile.Close()

	log.Println(" --------- Started")
	doCheckStructure := flag.Bool("check", true, "Check directory structure")
	doDownloadImages := flag.Bool("download", false, "Download images")
	doCreateGif := flag.Bool("gif", false, "Generate gif file")
	doServe := flag.Bool("serve", true, "Start the webserver")
	doPeriodicChecking := flag.Bool("periodicdownload", true, "Do hourly image check")
	gifDir := flag.String("gifdir", "gifs", "Directory for generated gifs")
	imagesDir := flag.String("imagesdir", "images", "Directory for downloaded images")
	frameTime := flag.Int("frametime", 50, "Delay between frames in 10ms")
	gifFileName := flag.String("gifname", "anim.gif", "Name of the produced gif file")
	address := flag.String("address", "0.0.0.0", "Local address to bind to")
	port := flag.Int("port", 8080, "Local port to bind to")

	flag.Parse()

	if *doCheckStructure {
		log.Println("Checking directory structure")
		ok := EnsureDirectoryStructure(*imagesDir, *gifDir)
		if !ok {
			log.Fatal("Failed to create directory structure")
		}
	}
	if *doDownloadImages {
		log.Println("Downloading images")
		DownloadImages(*imagesDir)
	}
	if *doCreateGif {
		log.Println("Creating gif")
		err := CreateGif(*frameTime, *imagesDir, path.Join(*gifDir, *gifFileName))
		if err != nil {
			log.Println(err)
		}
	}
	if *doPeriodicChecking {
		t := time.NewTicker(CHECKING_PERIOD)
		go func() {
			for range t.C {
				log.Println("Downloading images, periodic")
				downloaded := DownloadImages(*imagesDir)
				if downloaded > 0 {
					log.Println("Creating gif because image was downloaded")
					CreateGif(*frameTime, *imagesDir, path.Join(*gifDir, *gifFileName))
				}
			}
		}()
		defer t.Stop()
	}
	if *doServe {
		log.Println(" --------- Starting webserver")
		StartServer(*address, *port)
	}

	log.Println(" --------- Finished")
}
