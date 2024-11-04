package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var IMAP_SERVER string
var IMAP_PORT int
var IMAP_USER string
var IMAP_PASS string
var SKIP_FOLDERS []string
var IMAP_DEBUG bool
var SPECIFIC_FOLDER string

func main() {
	//Welcome print
	fmt.Println("###########################################")
	fmt.Println("\tWelcome to CleanInbox!")
	fmt.Println("")
	fmt.Printf("This program will help you clean your inbox\nor a specific folder by scanning for and\ndeleting duplicate emails.\n")
	fmt.Println("###########################################")

	loadENV()

	// Call the appropriate function based on user input
	keepAsking := true
	for keepAsking {
		fmt.Println("Do you want to scan or delete?")
		fmt.Print("Enter 'scan', 'scan-folder', 'delete', 'delete-folder' or 'exit': ")

		// Read user input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToLower(scanner.Text())

		switch input {
		case "scan":
			keepAsking = false
			scan()
		case "delete":
			keepAsking = false
			deleteAll()
		case "scan-folder":
			keepAsking = false
			listFolders()
			fmt.Print("Enter the folder name to scan (case sensitive): ")
			scanner.Scan()
			SPECIFIC_FOLDER = scanner.Text()
			scan()
		case "delete-folder":
			keepAsking = false
			listFolders()
			fmt.Print("Enter the folder name to delete from (case sensitive): ")
			scanner.Scan()
			SPECIFIC_FOLDER = scanner.Text()
			deleteAll()
		case "exit":
			keepAsking = false
		default:
			fmt.Println("Invalid choice. Please enter 'scan', 'delete', 'scan-folder', 'delete-folder' or 'exit'.")
		}
	}
}

func loadENV() {
	//Check if .env exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("No .env file found. Please create one with the following variables: IMAP_SERVER, IMAP_PORT, IMAP_USER, IMAP_PASS, SKIP_FOLDERS")

		//Copy .env.example to .env with os
		err2 := copyFile(".env.example", ".env")
		if err2 != nil {
			fmt.Printf("Error copying .env.example to .env: %v\n", err2)
		}

		os.Exit(1)
	}

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Load environment variables, strings
	IMAP_SERVER = os.Getenv("IMAP_SERVER")
	IMAP_USER = os.Getenv("IMAP_USER")
	IMAP_PASS = os.Getenv("IMAP_PASS")

	// Load environment variables, conversion needed
	IMAP_PORT_STRING := os.Getenv("IMAP_PORT")
	IMAP_DEBUG_STRING := os.Getenv("IMAP_DEBUG")

	// Load folders to skip
	SKIP_FOLDERS = strings.Split(os.Getenv("SKIP_FOLDERS"), ",")

	// Validate environment variables
	if IMAP_SERVER == "" || IMAP_PORT_STRING == "" || IMAP_USER == "" || IMAP_PASS == "" {
		fmt.Println("Please set the following environment variables: IMAP_SERVER, IMAP_PORT, IMAP_USER, IMAP_PASS")
		os.Exit(1)
	}

	// Convert IMAP_PORT to int
	IMAP_PORT, err = strconv.Atoi(IMAP_PORT_STRING)
	if err != nil {
		fmt.Println("IMAP_PORT must be valid port number between 1 and 65535, e.g. 993")
		os.Exit(1)
	}
	if IMAP_PORT < 1 || IMAP_PORT > 65535 {
		fmt.Println("IMAP_PORT must be valid port number between 1 and 65535, e.g. 993")
		os.Exit(1)
	}

	// Convert IMAP_DEBUG to bool
	IMAP_DEBUG, err = strconv.ParseBool(IMAP_DEBUG_STRING)
	if err != nil {
		fmt.Println("IMAP_DEBUG must be a boolean value, e.g. true or false. Default is false.")
		IMAP_DEBUG = false
	}

	// Print environment variables
	fmt.Println("IMAP_SERVER:", IMAP_SERVER)
	fmt.Println("IMAP_PORT:", IMAP_PORT)
	fmt.Println("IMAP_USER:", IMAP_USER)

	// Print folders to skip
	fmt.Println("SKIP_FOLDERS:")
	for _, folder := range SKIP_FOLDERS {
		fmt.Println("\t", folder)
	}

	// Print message
	fmt.Println("Press 'Ctrl + C' to exit at any time.")
}
