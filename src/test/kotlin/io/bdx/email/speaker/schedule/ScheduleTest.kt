package io.bdx.email.speaker.schedule

import com.fasterxml.jackson.databind.ObjectMapper
import io.bdx.email.speaker.event.EventTest
import org.assertj.core.api.Assertions.assertThat
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.junit.jupiter.SpringExtension

@SpringBootTest
@ExtendWith(SpringExtension::class)
class ScheduleTest {

    @Autowired
    private lateinit var objectMapper: ObjectMapper

    private val json = EventTest::class.java.getResource("/schedule.json").readText()

    @Test
    fun `it should load schedule from JSON`() {
        val schedule = Schedule.from(json, objectMapper)

        assertThat(schedule.sessions).hasSize(2)
    }
}