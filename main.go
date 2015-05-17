package main

import (
	"flag"
	_ "image/png"
	"log"
	"path"
)

func main() {
	log.Println(" --------- Started")
	doCheckStructure := flag.Bool("check", false, "Check directory structure")
	doDownloadImages := flag.Bool("download", false, "Download images")
	doCreateGif := flag.Bool("gif", false, "Generate gif file")
	gifDir := flag.String("gifdir", "gifs", "Directory for generated gifs")
	imagesDir := flag.String("imagesdir", "images", "Directory for downloaded images")
	frameTime := flag.Int("frametime", 50, "Delay between frames in 10ms")
	gifFileName := flag.String("gifName", "anim.gif", "Name of the produced gif file")

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
	log.Println(" --------- Finished")
}
