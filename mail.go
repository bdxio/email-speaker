package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/wneessen/go-mail"
)

type Data struct {
	Firstname       string
	Lastname        string
	TalkID          string
	Title           string
	Start           string
	End             string
	Room            string
	Backup          bool
	Quickie         bool
	Accommodation   bool
	OpenFeedbackURL string
	Female          bool
}

func CreateTemplate(path string) (*template.Template, error) {
	html, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return html, err
}

const openfeedBackURL = "https://openfeedback.io/r46KviPgLYMQfQnFpaGS/2022-12-02"

func CreateMail(subject string, invites bool, speaker Information, talk Talk, html *template.Template) (*mail.Msg, error) {
	log.Printf("Creating email for speaker %q with talk %q\n", speaker.UID, talk.Title)

	msg := mail.NewMsg()
	if err := msg.FromFormat("Team BDX I/O", "team@bdxio.fr"); err != nil {
		return nil, err
	}
	if err := msg.AddToFormat(fmt.Sprintf("%s %s", speaker.Firstname, speaker.Lastname), speaker.Email); err != nil {
		return nil, err
	}
	if err := msg.AddBcc("team@bdxio.fr"); err != nil {
		return nil, err
	}

	msg.Subject(subject)
	data := Data{
		Firstname:       speaker.Firstname,
		Lastname:        speaker.Lastname,
		TalkID:          talk.ID,
		Title:           talk.Title,
		Start:           talk.StartTime.Format("15h04"),
		End:             talk.EndTime.Format("15h04"),
		Room:            talk.Room,
		Backup:          talk.Backup,
		Quickie:         talk.IsQuickie(),
		Accommodation:   speaker.Accommodation,
		OpenFeedbackURL: openfeedBackURL,
		Female:          speaker.IsFemale(),
	}
	if err := msg.SetBodyHTMLTemplate(html, data); err != nil {
		log.Fatalf("failed to set HTML template as HTML body: %s", err)
	}

	if !invites {
		return msg, nil
	}

	// Might not be needed with the next tickets platform.
	//msg.AttachFile(getTicketFile(speaker))

	qr, err := generateQRCode(talk.ID)
	if err != nil {
		return nil, err
	}

	if !talk.Backup {
		msg.AttachReader("qr.png", qr)
		msg.AttachReader("talk.ics", generateTalkInvite(talk))
	}
	msg.AttachReader("dinner.ics", generateDinnerInvite())
	msg.AttachReader("after.ics", generateAfterInvite())

	return msg, nil
}

func getTicketFile(speaker Information) string {
	firstname := normalize(speaker.Firstname)
	lastname := normalize(speaker.Lastname)
	ticket := fmt.Sprintf("data/tickets/billet %s - %s %s.pdf", speaker.Ticket, firstname, lastname)
	if _, err := os.Stat(ticket); os.IsNotExist(err) {
		log.Fatalf("file %s does not exist", ticket)
	}
	return ticket
}

func normalize(s string) string {
	ns := strings.ToLower(s)
	ns = strings.ReplaceAll(ns, "é", "e")
	ns = strings.ReplaceAll(ns, "è", "e")
	ns = strings.ReplaceAll(ns, "î", "i")
	ns = strings.ReplaceAll(ns, "ï", "i")
	ns = strings.ReplaceAll(ns, "ô", "o")
	ns = strings.ReplaceAll(ns, "ç", "c")
	ns = strings.ReplaceAll(ns, "m. ", "")

	return ns
}

func generateQRCode(talkID string) (io.Reader, error) {
	data, err := qrcode.Encode(fmt.Sprintf("%s/%s", openfeedBackURL, talkID), qrcode.Highest, 768)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

const (
	icsBegin = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//bobbin v0.1//NONSGML iCal Writer//EN
CALSCALE:GREGORIAN
METHOD:PUBLISH`

	icsEnd = "END:VCALENDAR"
)

func generateTalkInvite(talk Talk) io.Reader {
	var b strings.Builder
	b.WriteString(icsBegin)
	b.WriteRune('\n')
	b.WriteString(generateEvent(
		talk.Title,
		fmt.Sprintf("%s, Palais des Congrès de Bordeaux, Av. Jean Gabriel Domergue, 33300 Bordeaux", talk.Room),
		44.8882631,
		-0.5694203,
		talk.StartTime,
		talk.EndTime,
	))
	b.WriteString(icsEnd)

	return strings.NewReader(b.String())
}

func generateDinnerInvite() io.Reader {
	var b strings.Builder
	b.WriteString(icsBegin)
	b.WriteRune('\n')
	b.WriteString(generateEvent(
		"Soirée Speakers",
		"La Brasserie Bordelaise, 50 Rue Saint-Rémi, 33000 Bordeaux",
		44.8411861,
		-0.5730474,
		time.Date(2022, 12, 1, 18, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 1, 22, 0, 0, 0, time.UTC),
	))
	b.WriteString(icsEnd)

	return strings.NewReader(b.String())
}

func generateAfterInvite() io.Reader {
	var b strings.Builder
	b.WriteString(icsBegin)
	b.WriteRune('\n')
	b.WriteString(generateEvent(
		"After BDX I/O",
		"La Cervoiserie, 10 Quai Lawton, 33300 Bordeaux",
		44.8411804,
		-0.5993119,
		time.Date(2022, 12, 2, 17, 30, 0, 0, time.UTC),
		time.Date(2022, 12, 2, 22, 0, 0, 0, time.UTC),
	))
	b.WriteString(icsEnd)

	return strings.NewReader(b.String())
}

func generateEvent(name, location string, lat, lon float64, start, end time.Time) string {
	var b strings.Builder
	b.WriteString("BEGIN:VEVENT\n")
	b.WriteString("DTSTART:")
	b.WriteString(start.Format("20060102T150405Z07"))
	b.WriteRune('\n')
	b.WriteString("DTEND:")
	b.WriteString(end.Format("20060102T150405Z07"))
	b.WriteRune('\n')
	b.WriteString("DTSTAMP:")
	b.WriteString(time.Now().Format("20060102T150405Z07"))
	b.WriteRune('\n')
	b.WriteString("UID:")
	b.WriteString(uuid.NewString())
	b.WriteRune('\n')
	b.WriteString("CREATED:")
	b.WriteString(time.Now().Format("20060102T150405Z07"))
	b.WriteRune('\n')
	b.WriteString("DESCRIPTION:")
	b.WriteString(name)
	b.WriteRune('\n')
	b.WriteString("LAST-MODIFIED:")
	b.WriteString(time.Now().Format("20060102T150405Z07"))
	b.WriteRune('\n')
	b.WriteString("SEQUENCE:0\n")
	b.WriteString("STATUS:CONFIRMED\n")
	b.WriteString("SUMMARY:")
	b.WriteString(name)
	b.WriteRune('\n')
	b.WriteString("ORGANIZER;CN=Team BDX I/O:MAILTO:team@bdxio.fr\n")
	b.WriteString("LOCATION:")
	b.WriteString(location)
	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("GEO:%f;%f\n", lat, lon))
	b.WriteString("TRANSP:OPAQUE\n")
	b.WriteString("END:VEVENT\n")

	return b.String()
}
