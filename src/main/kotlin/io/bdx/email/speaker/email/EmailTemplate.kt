package io.bdx.email.speaker.email

import com.github.mustachejava.Mustache
import com.github.mustachejava.MustacheFactory
import io.bdx.email.speaker.EmailProperties
import io.bdx.email.speaker.qrcode.QrCodeGenerator
import io.bdx.email.speaker.qrcode.generateQRCode
import io.bdx.email.speaker.schedule.Session
import io.bdx.email.speaker.schedule.SpeakerSchedule
import org.springframework.stereotype.Component
import java.io.File
import java.io.StringReader
import java.io.StringWriter
import java.time.Instant
import java.time.LocalDateTime
import java.time.ZoneId
import java.time.format.DateTimeFormatter

private val DATE_FORMATTER = DateTimeFormatter.ofPattern("d MMMM uuuu")
private val TIME_FORMATTER = DateTimeFormatter.ofPattern("HH'h'mm")
private val OPEN_FEEDBACK_URL_FORMATTER = DateTimeFormatter.ofPattern("uuuu-MM-dd")

@Component
class EmailTemplate(
    private val emailProperties: EmailProperties,
    private val mustacheFactory: MustacheFactory,
    private val qrCodeGenerator: QrCodeGenerator
) {

    fun generateEmails(speakerSchedules: List<SpeakerSchedule>): List<Email> {
        val template = File(emailProperties.templatePath).readText()
        val mustache = mustacheFactory.compile(StringReader(template), "email")
        return speakerSchedules.map { generateEmail(it, mustache) }
    }

    private fun generateEmail(speakerSchedule: SpeakerSchedule, mustache: Mustache): Email {
        val writer = StringWriter()
        mustache.execute(writer, speakerSchedule.scopes()).flush()
        val qrCodes = generateQRCodes(speakerSchedule.sessions)
        val recipient = if (emailProperties.to.isBlank()) speakerSchedule.speaker.email else emailProperties.to
        return Email(emailProperties.from, emailProperties.to, emailProperties.subject, writer.toString(), qrCodes)
    }

    private fun generateQRCodes(sessions: List<Session>) = sessions
        .map {
            Pair(
                "${it.title} QR code.png",
                emailProperties.openFeedbackUrl +
                        "/" +
                        OPEN_FEEDBACK_URL_FORMATTER.format(it.startTime.toLocalDateTime()) +
                        "/${it.id}"
            )
        }
        .map { EmailAttachment(it.first, qrCodeGenerator.generate(it.second)) }
}

private fun SpeakerSchedule.scopes() = mapOf(
    "displayName" to speaker.displayName,
    "sessions" to sessions.scopes()
)

private fun List<Session>.scopes() = map { it.scopes() }

private fun Session.scopes() = mapOf(
    "title" to title,
    "startDate" to DATE_FORMATTER.format(startTime.toLocalDateTime()),
    "startTime" to TIME_FORMATTER.format(startTime.toLocalDateTime()),
    "endTime" to TIME_FORMATTER.format(endTime.toLocalDateTime()),
    "room" to trackTitle
)

private fun Instant.toLocalDateTime() = LocalDateTime.ofInstant(this, ZoneId.of("Europe/Paris"))
