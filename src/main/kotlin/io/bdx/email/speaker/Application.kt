package io.bdx.email.speaker

import com.fasterxml.jackson.databind.DeserializationFeature
import com.fasterxml.jackson.databind.MapperFeature
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule
import com.fasterxml.jackson.module.kotlin.registerKotlinModule
import com.github.mustachejava.DefaultMustacheFactory
import com.github.mustachejava.MustacheFactory
import io.bdx.email.speaker.email.EmailSender
import io.bdx.email.speaker.email.EmailTemplate
import io.bdx.email.speaker.schedule.SpeakerSchedulesRetriever
import org.springframework.boot.CommandLineRunner
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.ConstructorBinding
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.mail.javamail.JavaMailSender
import org.springframework.mail.javamail.JavaMailSenderImpl

@SpringBootApplication
@EnableConfigurationProperties
class Application(
    private val speakerSchedulesRetriever: SpeakerSchedulesRetriever,
    private val emailTemplate: EmailTemplate,
    private val emailSender: EmailSender
) : CommandLineRunner {

    override fun run(vararg args: String?) {
        val speakerSchedules = speakerSchedulesRetriever.getSpeakerSchedules()
        val emails = emailTemplate.generateEmails(speakerSchedules)
        emailSender.sendEmails(emails.take(5))
    }
}

fun main(args: Array<String>) {
    runApplication<Application>(*args)
}

@Configuration
class Configuration(private val emailProperties: EmailProperties) {
    @Bean
    fun objectMapper(): ObjectMapper = ObjectMapper()
        .registerKotlinModule()
        .registerModule(JavaTimeModule())
        .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
        .enable(MapperFeature.ACCEPT_CASE_INSENSITIVE_ENUMS)

    @Bean
    fun mustacheFactory(): MustacheFactory = DefaultMustacheFactory()

    @Bean
    fun javaMailSender(): JavaMailSender {
        System.setProperty("mail.mime.splitlongparameters", "false")
        val sender = JavaMailSenderImpl()
        sender.host = emailProperties.host
        sender.port = emailProperties.port
        sender.username = emailProperties.username
        sender.password = emailProperties.password
        with(sender.javaMailProperties) {
            put("mail.transport.protocol", "smtp")
            put("mail.smtp.auth", "true")
            put("mail.smtp.starttls.enable", "true")
        }

        return sender
    }
}

@ConstructorBinding
@ConfigurationProperties(prefix = "input")
data class InputProperties(val eventPath: String, val scheduleUrl: String)

@ConstructorBinding
@ConfigurationProperties(prefix = "email")
data class EmailProperties(
    val templatePath: String,
    val host: String,
    val port: Int,
    val username: String,
    val password: String,
    val from: String,
    val to: String,
    val subject: String,
    val openFeedbackUrl: String
)