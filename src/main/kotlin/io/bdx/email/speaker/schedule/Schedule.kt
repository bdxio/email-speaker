package io.bdx.email.speaker.schedule

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.readValue
import io.bdx.email.speaker.event.Talk
import java.time.Instant

data class Schedule(val sessions: Map<String, Session>) {
    companion object {
        fun from(json: String, objectMapper: ObjectMapper) = objectMapper.readValue<Schedule>(json)
    }
}

data class Session(
    val id: Int,
    val title: String,
    val startTime: Instant,
    val endTime: Instant,
    val trackTitle: String
)

data class ScheduledTalk(val talk: Talk, val session: Session)
