package com.example.sunriseappnew.model

import java.util.Calendar

internal class Alarm(sunrise: String) {

    private var hour: Int = 0
    private var minute: Int = 0

    init {
        stringToTime(sunrise)
    }

    override fun toString(): String {
        var time = ""
        if (minute != 0 || hour != 0) {
            time = if (minute < 10) {
                "$hour:0$minute"
            } else {
                "$hour:$minute"
            }
        }
        return "$time AM"
    }

    fun timeToMilli(currentHour: Int, currentMinute: Int): Long {
        val unix: Long
        val calendar = Calendar.getInstance()
        val current =
            (((calendar[Calendar.HOUR_OF_DAY] * 60) + calendar[Calendar.MINUTE]) * 60000).toLong()
        val milli = (currentHour.toLong() * 60 + currentMinute.toLong()) * 60000
        unix = if (milli > current) {
            milli - current
        } else {
            ((24 * 60) * 60000) - (current - milli)
        }
        return unix
    }

    private fun stringToTime(officialSunrise: String) {
        val arr = officialSunrise.toCharArray()
        for (i in 0..1) {
            this.minute += if (i == 0) {
                Character.getNumericValue(arr[i]) * 10
            } else {
                Character.getNumericValue(arr[i])
            }
        }
        for (i in 3..4) {
            this.hour += if (i == 3) {
                Character.getNumericValue(arr[i]) * 10
            } else {
                Character.getNumericValue(arr[i])
            }
        }
    }
}