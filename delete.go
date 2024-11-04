package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func deleteAll() {
	im := emailClient()
	defer im.Close()

	rootDomainsList := loadFile()

	// Folders now contains a string slice of all the folder names on the connection
	folders, err := im.GetFolders()
	check(err)

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
		err = im.SelectFolder(f)
		check(err)

		// This function implements the IMAP UID search, returning a slice of ints
		// Sending "ALL" runs the command "UID SEARCH ALL"
		// You can enter things like "*:1" to get the first UID, or "999999999:*"
		// to get the last (unless you actually have more than that many emails)
		// You can check out https://tools.ietf.org/html/rfc3501#section-6.4.4 for more
		uids, err := im.GetUIDs("ALL")
		check(err)

		// GetEmails takes a list of ints as UIDs, and returns new Email objects.
		// If an email for a given UID cannot be found, there's an error parsing its body,
		// or the email addresses are malformed (like, missing parts of the address), then it is skipped
		// If an email is found, then an imap.Email struct slice is returned with the information from the email.
		emails, err := im.GetEmails(uids...)
		check(err)

		if len(emails) != 0 {
			// Print a summary of one of the emails
			// (note: the emails may not be returned in any particular order)
			// fmt.Print(emails[0])

			// Count emails per sender
			for uid, addresses := range emails {
				delete := false
				for emailAddress, _ := range addresses.From {
					// fmt.Printf("%s, %s\n", email, name)
					for listDomain, _ := range rootDomainsList {
						//Check if its an email or a domain
						if strings.Contains(listDomain, "@") {
							if listDomain == emailAddress {
								delete = true
							}
						} else {
							if listDomain == extractRootDomain(emailAddress) {
								delete = true
							}
						}
					}
				}

				if delete {
					// Mark the matching emails for deletion
					// item := seqSet.Items[0]
					// set := &imap.SeqSet{}
					// set.AddNum(item)

					// // Set the \Deleted flag
					// items := make(map[uint32]imap.FlagSet)
					// items[item] = []string{imap.DeletedFlag}
					// if err := im.Store(set, items, nil); err != nil {
					// 	log.Fatalf("Failed to mark emails for deletion: %v", err)
					// 	return
					// }

					im.MoveEmail(uid, "Trash")
				}

			}

		}
	}
}

func loadFile() map[string]int {
	// Read root domains and counts from delete.txt
	file, err := os.Open("delete.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a map to store root domains and counts
	rootDomainCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			rootDomain := strings.TrimSpace(parts[0])
			count := strings.TrimSpace(parts[1])
			// Convert count to integer
			if countInt, err := strconv.Atoi(count); err == nil {
				rootDomainCounts[rootDomain] = countInt
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return rootDomainCounts
}
