package io.bdx.email.speaker.event

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.readValue
import io.bdx.email.speaker.schedule.Session
import java.time.Instant

data class Event(
    val name: String,
    val conferenceDates: ConferenceDates,
    val talks: List<Talk>,
    val speakers: List<Speaker>
) {
    fun confirmedTalks() = talks.filter { it.isConfirmed() }

    fun confirmedSpeakers() =
        confirmedTalks().flatMap { it.speakers }
            .distinct()
            .map { uid -> speakers.first { it.uid == uid } }

    fun confirmedSpeakersWithConfirmedTalks() =
        confirmedSpeakers()
            .map { speaker ->
                SpeakerWithTalks(speaker, confirmedTalks().filter { it.speakers.contains(speaker.uid) })
            }

    companion object {
        fun from(json: String, objectMapper: ObjectMapper) = objectMapper.readValue<Event>(json)
    }
}

data class ConferenceDates(val start: Instant, val end: Instant)

typealias UID = String

data class Talk(val title: String, val state: State, val speakers: List<UID>) {
    enum class State { SUBMITTED, REJECTED, CONFIRMED, DECLINED }

    fun isConfirmed() = state == State.CONFIRMED

    fun session(sessions: Collection<Session>): Session =
        sessions.firstOrNull { it.title.trim() == title.trim() }
                ?: throw IllegalArgumentException("no session found for talk $title")
}

data class Speaker(val uid: UID, val displayName: String?, val email: String?)

data class SpeakerWithTalks(val speaker: Speaker, val talks: List<Talk>)