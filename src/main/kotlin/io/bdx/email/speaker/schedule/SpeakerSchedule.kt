package io.bdx.email.speaker.schedule

import com.fasterxml.jackson.databind.ObjectMapper
import com.github.kittinunf.fuel.Fuel
import com.github.kittinunf.fuel.jackson.responseObject
import com.github.kittinunf.result.Result
import io.bdx.email.speaker.InputProperties
import io.bdx.email.speaker.event.Event
import io.bdx.email.speaker.event.Speaker
import org.springframework.stereotype.Component
import java.io.File

data class SpeakerSchedule(val speaker: Speaker, val sessions: List<Session>)

@Component
class SpeakerSchedulesRetriever(private val props: InputProperties, private val objectMapper: ObjectMapper) {

    fun getSpeakerSchedules(): List<SpeakerSchedule> {
        val event = Event.from(File(props.eventPath).readText(), objectMapper)

        val (_, _, result) =
            Fuel.get(props.scheduleUrl)
                .responseObject<Schedule>(objectMapper)
        val schedule = when (result) {
            is Result.Failure -> throw result.getException()
            is Result.Success -> result.value
        }

        return event.confirmedSpeakersWithConfirmedTalks()
            .map { speakerWithTalks ->
                SpeakerSchedule(
                    speakerWithTalks.speaker,
                    speakerWithTalks.talks.map { it.session(schedule.sessions.values) })
            }
    }
}
