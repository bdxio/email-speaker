package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wneessen/go-mail"
)

func main() {
	var (
		subject  string
		template string
		invites  bool
		url      string
		port     int
		username string
		password string
		dryRun   bool
	)
	flag.StringVar(&subject, "subject", "", "subject of the email")
	flag.StringVar(&template, "template", "", "path to the template to use for the email")
	flag.BoolVar(&invites, "invites", false, "to include talk, dinner and after invites as well as the talk feedback QR code in the email")
	flag.StringVar(&url, "url", "", "URL of the email server")
	flag.IntVar(&port, "port", 587, "port of the email server")
	flag.StringVar(&username, "username", "", "username of the email server")
	flag.StringVar(&password, "password", "", "password of the email server")
	flag.BoolVar(&dryRun, "dry-run", true, "set to false to really send the emails")
	flag.Parse()

	if subject == "" {
		flag.Usage()
		log.Fatalf("the subject of the email is required")
	}

	if template == "" {
		flag.Usage()
		log.Fatalf("the path to the template to use for the email is required")
	}

	if dryRun == false && (url == "" || username == "" || password == "") {
		flag.Usage()
		log.Fatalf("invalid argument")
	}

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatalf("Unable to load time location for Europe/Paris: %v", err)
	}

	talks, err := ReadTalks(loc)
	if err != nil {
		log.Fatalf("Unable to read talks from JSON: %v", err)
	}

	speakers, err := ReadSpeakers()
	if err != nil {
		log.Fatalf("Unable to read speakers from CSV: %v", err)
	}

	// Add missing talks
	// Closing keynote
	closingKeynote := Talk{
		Title:     "✏️ La pensée visuelle a changé ma vie",
		StartTime: time.Date(2022, 12, 2, 17, 20, 0, 0, loc),
		EndTime:   time.Date(2022, 12, 2, 18, 5, 0, 0, loc),
		Room:      "Amphi A",
	}
	// Add talk to each speaker
	talks["47988259"] = []Talk{closingKeynote} // Horacio Gonzales
	talks["47988262"] = []Talk{closingKeynote} // Pierre Tibule
	talks["47988263"] = []Talk{closingKeynote} // Aurélie Vache

	html, err := CreateTemplate(template)
	if err != nil {
		log.Fatalf("Unable to create HTML template: %v", err)
	}

	msgs := make([]*mail.Msg, 0, len(talks))
	for speakerID, speakerTalks := range talks {
		if len(speakerTalks) > 1 {
			log.Fatalf("Cannot handle more than one talk for speaker with ID %s", speakerID)
		}

		speaker, ok := speakers[speakerID]
		if !ok {
			log.Fatalf("Unable to find speaker with ID %q for talk %s", speakerID, speakerTalks[0].Title)
		}

		talk := speakerTalks[0]

		msg, err := CreateMail(subject, invites, speaker, talk, html)
		if err != nil {
			log.Fatalf("Unable to create email for speaker %v: %v", speaker, err)
		}

		msgs = append(msgs, msg)
	}

	if dryRun {
		log.Printf("Writing %d emails to output folder", len(msgs))
		if err := os.Mkdir("output", 0o644); err != nil && !errors.Is(err, os.ErrExist) {
			log.Fatalf("Could not create output folder: %v", err)
		}
		for _, msg := range msgs {
			if err := msg.WriteToFile(fmt.Sprintf("output/%s.txt", msg.GetTo()[0])); err != nil {
				log.Fatalf("Could not write email to %s to disk: %v", msg.GetTo(), err)
			}
		}
		return
	}

	c, err := mail.NewClient(url,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)
	if err != nil {
		log.Fatalf("Unable to create client for mail server: %v", err)
	}
	defer c.Close()

	log.Printf("Sending %d emails...", len(msgs))
	if err := c.DialAndSend(msgs...); err != nil {
		log.Fatalf("Unable to send emails: %v", err)
	}

	log.Println("All done.")
}
