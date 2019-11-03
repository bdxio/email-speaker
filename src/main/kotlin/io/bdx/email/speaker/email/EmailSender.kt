package io.bdx.email.speaker.email

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.mail.javamail.JavaMailSender
import org.springframework.mail.javamail.MimeMessageHelper
import org.springframework.stereotype.Component
import javax.mail.internet.MimeUtility

@Component
class EmailSender {
    @Autowired
    private lateinit var sender: JavaMailSender

    fun sendEmails(emails: List<Email>) = emails.forEach { sendEmail(it) }

    private fun sendEmail(email: Email) {
        val message = sender.createMimeMessage()
        val helper = MimeMessageHelper(message, true, Charsets.UTF_8.name())
        helper.setFrom(email.from)
        helper.setTo(email.to)
        helper.setSubject(email.subject)
        helper.setText(email.text, true)
        email.attachments.forEach { helper.addAttachment(MimeUtility.encodeWord(it.filename), it.file) }

        sender.send(message)
    }
}
