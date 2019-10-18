package main

type Ticket struct {
	ValidDate    string    `json:"validDate"`
	PdfDesign    string    `json:"pdfDesign"`
	PdfWidth     float64   `json:"pdfWidth"`
	PdfHeight    float64   `json:"pdfHeight"`
	XCoordinates []float64 `json:"x_coordinates"`
	YCoordinates []float64 `json:"y_Coordinates"`
}

// Configuration will contain all config data

type Configuration struct {
	Ticket             Ticket `json:"ticket"`
	CsvFile            string `json:"csv_file"`
	Seed 			   int64 `json:"seed"`
	Password			   string`json:"seed"`
	DBName				string `json:"db_name"`
}

//Student to contain all info needed for ticket creation

type Student struct {
	Name       string // Student name
	Dni        string // Student ID (only numbers)
	NumTickets int    // Number of tickets that will be created for this student
	Keys       []string
}