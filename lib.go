package main

import (
	"io"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// List of public email providers not to group
var PUBLIC_DOMAINS = []string{
	"gmail.com",
	"yahoo.com",
	"hotmail.com",
	"outlook.com",
	"icloud.com",
	"live.com",
	"me.com",
	"msn.com",
	"ymail.com",
	"rocketmail.com",
	"aim.com",
	"zoho.com",
	"protonmail.com",
	"tutanota.com",
	"mail.com",
	"yandex.com",
	"inbox.com",
	"gmx.com",
	"fastmail.com",
	"lavabit.com",
	"runbox.com",
	"posteo.de",
	"kolabnow.com",
	"disroot.org",
	"riseup.net",
	"autistici.org",
	"tuta.io",
	"mailbox.org",
	"countermail.com",
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func extractRootDomain(email string) string {
	// Split the email address by '@'
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}

	etld1, err := publicsuffix.EffectiveTLDPlusOne(parts[1])
	suffix, icann := publicsuffix.PublicSuffix(strings.ToLower(parts[1]))
	if err != nil && !icann && suffix == parts[1] {
		etld1 = parts[1]
		err = nil
	}
	if err != nil {
		return ""
	}

	return etld1
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the content from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush the file to disk
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
