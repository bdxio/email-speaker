package io.bdx.email.speaker.qrcode

import net.glxn.qrgen.javase.QRCode
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.ConstructorBinding
import org.springframework.stereotype.Component
import java.io.File

@Component
class QrCodeGenerator(private val qrCodeProperties: QrCodeProperties) {
    fun generate(text: String): File =
        QRCode.from(text).withSize(qrCodeProperties.width, qrCodeProperties.height).file()
}

@ConstructorBinding
@ConfigurationProperties(prefix = "qr-code")
data class QrCodeProperties(val width: Int, val height: Int)