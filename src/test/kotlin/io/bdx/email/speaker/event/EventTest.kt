package io.bdx.email.speaker.event

import com.fasterxml.jackson.databind.ObjectMapper
import io.bdx.email.speaker.schedule.Session
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.junit.jupiter.SpringExtension
import java.time.Instant

@SpringBootTest
@ExtendWith(SpringExtension::class)
class EventTest {

    @Autowired
    private lateinit var objectMapper: ObjectMapper

    private val json = EventTest::class.java.getResource("/event.json").readText()

    @Test
    fun `it should load event from JSON`() {
        val event = Event.from(json, objectMapper)

        assertThat(event.name).isEqualTo("BDX I/O")
        assertThat(event.talks).hasSize(5)
        assertThat(event.speakers).hasSize(1)
    }

    @Test
    fun `it should get confirmed talks`() {
        val event = Event.from(json, objectMapper)

        val confirmedTalks = event.confirmedTalks()

        assertThat(confirmedTalks).hasSize(2)
    }

    @Test
    fun `it should get confirmed speakers`() {
        val event = Event.from(json, objectMapper)

        val confirmedSpeakers = event.confirmedSpeakers()

        assertThat(confirmedSpeakers).hasSize(1)
    }

    @Test
    fun `it should get confirmed speakers with confirmed talks`() {
        val event = Event.from(json, objectMapper)

        val speakersWithTalks = event.confirmedSpeakersWithConfirmedTalks()

        assertThat(speakersWithTalks).hasSize(1)
        assertThat(speakersWithTalks.first().talks).hasSize(2)
        assertThat(speakersWithTalks.first().talks.first().title).isEqualTo("Une super présentation")
        assertThat(speakersWithTalks.first().talks.last().title).isEqualTo("Une autre super présentation")
    }

    @Test
    fun `it should get session for talk`() {
        val talk = Talk("Title", Talk.State.CONFIRMED, listOf("uid"))
        val sessions = listOf(
            Session(0, "Another Title", Instant.now(), Instant.now(), "Room 1"),
            Session(1, "Title", Instant.now(), Instant.now(), "Room 2")
        )

        val session = talk.session(sessions)

        assertThat(session.title).isEqualTo("Title")
        assertThat(session.trackTitle).isEqualTo("Room 2")
    }
}