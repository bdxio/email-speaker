package io.bdx.email.speaker.email

import java.io.File

data class Email(
    val from: String,
    val to: String,
    val subject: String,
    val text: String,
    val attachments: List<EmailAttachment>
)

data class EmailAttachment(val filename: String, val file: File)