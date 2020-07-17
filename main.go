package main

import (
	"flag"
	"log"

	"git.onebytedata.com/OneByteDataPlatform/go-dicom/media"
	"git.onebytedata.com/OneByteDataPlatform/go-dicom/services"
	"git.onebytedata.com/OneByteDataPlatform/one-byte-module/models"
)

var destination *models.Destination

func main() {
	media.InitDict()

	hostName := flag.String("host", "localhost", "Destination host name or IP")
	calledAE := flag.String("calledae", "DICOM_SCP", "AE of the destination")
	callingAE := flag.String("callingae", "DICOM_SCU", "AE of the client")
	port := flag.Int("port", 104, "Port of the destination system")

	studyUID := flag.String("studyuid", "", "Study UID to be added to request")

	destinationAE := flag.String("destinationae", "", "AE of the destination for a C-Move request")

	fileName := flag.String("file", "", "DICOM file to be sent")

	cecho := flag.Bool("cecho", false, "Send C-Echo to the destination")
	cfind := flag.Bool("cfind", false, "Send C-Find request to the destination")
	cmove := flag.Bool("cmove", false, "Send C-Move request to the destination")
	cstore := flag.Bool("cstore", false, "Sends a C-Store request to the destination")

	flag.Parse()

	destination = &models.Destination{
		Name:      *hostName,
		HostName:  *hostName,
		CalledAE:  *calledAE,
		CallingAE: *callingAE,
		Port:      *port,
		IsCFind:   true,
		IsCStore:  true,
		IsMWL:     true,
		IsTLS:     false,
	}

	if *cecho {
		scu := services.NewSCU(destination)
		err := scu.EchoSCU(30)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CEcho was successful")
	}
	if *cfind {
		query := media.DefaultCFindRequest()
		scu := services.NewSCU(destination)

		results := make([]media.DcmObj, 0)
		_, err := scu.FindSCU(query, &results, 30)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CFind was successful")
		log.Printf("Found %d results\n\n", len(results))
		for _, result := range results {
			log.Printf("Found study %s\n", result.GetString(0x0020, 0x000D))
			result.DumpTags()
		}
	}
	if *cmove {
		if *destinationAE == "" {
			log.Fatalln("destinationae is required for a C-Move")
		}
		if *studyUID == "" {
			log.Fatalln("studyuid is required for a C-Move")
		}
		query := media.DefaultCMoveRequest(*studyUID)
		scu := services.NewSCU(destination)
		_, err := scu.MoveSCU(*destinationAE, query, 30)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CMove was successful")
	}
	if *cstore {
		if *fileName == "" {
			log.Fatalln("file is required for a C-Store")
		}
		scu := services.NewSCU(destination)
		err := scu.StoreSCU(*fileName, 30)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
