package main

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/BrianLeishman/go-imap"
)

// SenderCount is a struct to store sender and count information
type SenderCount struct {
	Sender string
	Count  int
}

func listFolders() {
	im := emailClient()
	defer im.Close()

	// Folders now contains a string slice of all the folder names on the connection
	folders, err := im.GetFolders()
	check(err)

	fmt.Println("\tFolders:")
	for _, f := range folders {
		fmt.Printf("\t- %s\n", f)
	}
}

func scan() {
	im := emailClient()
	defer im.Close()

	// Create delete.txt file
	deleteFile, err := os.Create("delete.txt")
	if err != nil {
		panic(err)
	}
	defer deleteFile.Close()

	// Folders now contains a string slice of all the folder names on the connection
	folders, err := im.GetFolders()
	check(err)

	senderCount := make(map[string]int)

	// Now we can loop through those folders
	for _, f := range folders {
		if SPECIFIC_FOLDER != "" {
			if f != SPECIFIC_FOLDER {
				//Silent skip, they already know were targetting a folder
				continue
			}
		} else {
			// Skip some folders in env
			if slices.Contains(SKIP_FOLDERS, f) {
				fmt.Printf("Skipping folder %s\n", f)
				continue
			} else {
				fmt.Printf("Scanning Folder %s...\n", f)
			}
		}

		// And select each folder, one at a time.
		// Whichever folder is selected last, is the current active folder.
		// All following commands will be executing inside of this folder
		fmt.Println("\tSelecting Folder...")
		err = im.SelectFolder(f)
		check(err)

		// This function implements the IMAP UID search, returning a slice of ints
		// Sending "ALL" runs the command "UID SEARCH ALL"
		// You can enter things like "*:1" to get the first UID, or "999999999:*"
		// to get the last (unless you actually have more than that many emails)
		// You can check out https://tools.ietf.org/html/rfc3501#section-6.4.4 for more
		fmt.Println("\tGetting UIDs...")
		uids, err := im.GetUIDs("ALL")
		check(err)

		fmt.Printf("\t\tFound %d UIDs\n", len(uids))

		// GetEmails takes a list of ints as UIDs, and returns new Email objects.
		// If an email for a given UID cannot be found, there's an error parsing its body,
		// or the email addresses are malformed (like, missing parts of the address), then it is skipped
		// If an email is found, then an imap.Email struct slice is returned with the information from the email.
		fmt.Println("\tGetting Emails...")

		//Split the UIDs into chunks of 1000
		allEmails := make(map[int]*imap.Email)
		for i := 0; i < len(uids); i += 1000 {
			end := i + 1000
			if end > len(uids) {
				end = len(uids)
			}

			fmt.Printf("\t\tGetting emails %d to %d\n", i, end)
			emails, err := im.GetEmails(uids[i:end]...)
			check(err)

			fmt.Printf("\t\t\tFound %d emails\n", len(emails))
			for _, email := range emails {
				allEmails[email.UID] = email
			}
		}

		// emails, err := im.GetEmails(uids...)
		// check(err)

		fmt.Println("\tParsing Emails...")

		if len(allEmails) != 0 {
			// fmt.Printf("Found %d emails\n", len(emails))
			// Print a summary of one of the emails
			// (note: the emails may not be returned in any particular order)
			// fmt.Print(emails[0])

			// Count emails per sender
			for _, addresses := range allEmails {
				for email := range addresses.From {
					// fmt.Printf("%s, %s\n", email, name)
					rootDomain := extractRootDomain(email)

					//Check if its part of a public domain
					if slices.Contains(PUBLIC_DOMAINS, rootDomain) {
						senderCount[email]++
					} else {
						senderCount[rootDomain]++
					}
				}
			}

		}

		fmt.Println("\tDone.")
	}

	// os.Exit(1)

	fmt.Println("Done scanning, parsing results...")

	// Convert the map to a slice for sorting
	var senderCountSlice []SenderCount
	for sender, count := range senderCount {
		senderCountSlice = append(senderCountSlice, SenderCount{Sender: sender, Count: count})
	}

	// Sort the slice by count in ascending order
	sort.Slice(senderCountSlice, func(i, j int) bool {
		return senderCountSlice[i].Count < senderCountSlice[j].Count
	})

	// Print the count per sender
	fmt.Println("Emails per sender:")
	for _, sc := range senderCountSlice {
		if sc.Count > 1 {
			deleteFile.WriteString(fmt.Sprintf("%s: %d\n", sc.Sender, sc.Count))
			fmt.Printf("%s: %d\n", sc.Sender, sc.Count)
		}
	}
	deleteFile.Sync()

	fmt.Println("Results saved to delete.txt")
	fmt.Println("Remove any emails from delete.txt that you do not want to delete and rerun cleanInbox with 'delete' or 'delete-folder' option.")
	fmt.Println("Done.")
}
