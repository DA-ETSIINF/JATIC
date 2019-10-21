package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var students []Student
/// All names will be loaded from a file, probably csv

var configuration  = Configuration{}

func init() {
	LoadConfiguration("config.json")
	// validate configuration, specifically check for same number of X and Y Coordinates and number of fields to
	// autocomplete

	if len(configuration.Ticket.XCoordinates) != 3||
		len(configuration.Ticket.XCoordinates) != len(configuration.Ticket.YCoordinates) {
		log.Println(configuration.Seed)
		log.Fatal("Configuration error")
	}
	checkDatabase()
	rand.Seed(configuration.Seed)
}

func LoadConfiguration(file string) Configuration {

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	er := jsonParser.Decode(&configuration)
	if er != nil {
		log.Fatal(er)
	}
	return configuration
}

func main() {

	loadData()

	createKeys(students) // createKeys will globally modify existing data that we will pass on
	createPdf(students) // we only use already existing data and will not need it modified anywhere else
	insertData(students)


}


// loadData will try to load all data needed from given files
// it returns an array of students
func loadData() []Student {
	csvFile, _ := os.Open(configuration.CsvFile)
	r := csv.NewReader(bufio.NewReader(csvFile))

	for {
		record, err := r.Read()

		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			log.Println("Error loading data")
			log.Fatal(err)
		}
		numTickets, _  := strconv.Atoi(record[2])
		students = append(students, Student{
			Name: record[0],
			Dni: record[1],
			NumTickets: numTickets,

		})
	}
	return nil
}

/// Based on seed from config, userDni, number of tickets and a rand Int it will generate a SHA1 hash
// this will be used later on to create a the QR

func createKeys(students []Student) {
	for _, student := range students {
		for i := 0; i < student.NumTickets; i++ {
			seedString := student.Dni + configuration.Password + strconv.Itoa(rand.Int()) + strconv.Itoa(i)
			h := sha1.New()
			h.Write([]byte(seedString)) //generate hash  with the seedString
			bs := h.Sum(nil) // return hash as byte slice
			log.Print("Inserting key ticket num")
			log.Print(i)
			log.Print( base64.StdEncoding.EncodeToString(bs))
			// add it as a key for generating a ticket
			student.Keys = append(student.Keys, base64.StdEncoding.EncodeToString(bs))




		}
	}


}



// createPdf will create a folder (if it doesn't exist ) to store Tickets
// all Tickets created will have a name like Dni_NumTickets.pdf ie 55555555_01.pdf
// first it will create an  unique id for the Ticket
// second unique id will be saved
// third a QR  for that code will be created
// four pdf will be mounted
// name of pdf will be related to its unique id

func createPdf(students []Student)  {
	pdf :=  gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: gopdf.Rect{
		//Ticket size in points
		W: configuration.Ticket.PdfWidth,
		H: configuration.Ticket.PdfHeight,
	} })
	for _, student := range students {
		for ticknum := 0; ticknum < len(student.Keys); ticknum++  {

			pdf.AddPage()

			var err error

			err = pdf.AddTTFFont("DejaVu", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
			if err != nil {
				log.Println("Font error")
				log.Fatalf(err.Error())
				return
			}


			err = pdf.SetFont("DejaVu", "", 14)
			if err != nil {
				log.Print(err.Error())
			}

			template := pdf.ImportPage(configuration.Ticket.PdfDesign, 1, "/MediaBox" )

			pdf.UseImportedTemplate(template, 0, 0, 0, 0)

			for i := 0; i<len(configuration.Ticket.XCoordinates); i++ {
				/// Set coordinates for menu number
				pdf.SetX(configuration.Ticket.XCoordinates[i])
				pdf.SetY(configuration.Ticket.YCoordinates[i])

				// Find a better way to dinamically select data to write
				if i == 0{
					err = pdf.Cell(nil, student.Dni+strconv.Itoa(ticknum) )
					if err != nil {
						log.Print(err.Error())
					}
				}else if i==1 {
					err = pdf.Cell(nil, configuration.Ticket.ValidDate)
					if err != nil {
						log.Println("PDF num write error")
						log.Print(err.Error())
					}
				}else  {
					// last thing to add is a qr with SHA1 key

					err := qrcode.WriteFile(student.Keys[ticknum], qrcode.Highest, 180, "qr.png")
					err = pdf.Image("qr.png", configuration.Ticket.XCoordinates[i],
						configuration.Ticket.YCoordinates[i], nil)
					if err != nil {
						log.Println("Error creating qr")
						log.Print(err.Error())
					}
				}
			}
		}
	}
	//fileName := fmt.Sprintf("tickets/%s_%d.pdf", student.Dni, ticknum)

	err = pdf.WritePdf("tickets.pdf")
	if err != nil {
		log.Println("Error saving PDF")
		fmt.Print(err.Error())
	}
	err = pdf.Close()
	if err != nil {
		log.Println("Error saving PDF")
		fmt.Print(err.Error())
	}
}
