package main

import (
	"io"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	// TODO: adjust IMAP URL
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	// TODO: adjust mail and password
	if err = c.Login("max.mustermann@gmail.com", "1234"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err = <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}

	from := uint32(1)
	to := mbox.Messages
	// Get the last 10 messages
	if mbox.Messages > 9 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages)
	}()

	for msg := range messages {
		var section imap.BodySectionName

		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		// Print some info about the message
		header := mr.Header
		if date, err := header.Date(); err == nil {
			log.Println("Date:", date)
		}
		// if from, err := header.AddressList("From"); err == nil {
		// 	log.Println("From:", from)
		// }
		// if to, err := header.AddressList("To"); err == nil {
		// 	log.Println("To:", to)
		// }
		if subject, err := header.Subject(); err == nil {
			log.Println("Subject:", subject)
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				// b, _ := ioutil.ReadAll(p.Body)
				// log.Println("Got text: %v", string(b))
			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				log.Printf("Got attachment: %v", filename)
				// Create file with attachment name
				file, err := os.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				// using io.Copy instead of io.ReadAll to avoid insufficient memory issues
				size, err := io.Copy(file, p.Body)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Saved %v bytes into %v\n", size, filename)
			}
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}
