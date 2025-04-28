package com.example.sunriseappnew.model

import android.content.Context
import android.content.Intent
import android.os.Build
import android.provider.AlarmClock
import android.util.Log
import androidx.annotation.RequiresApi
import java.util.Calendar
import com.example.sunriseappnew.view.days

/**
 * Launches an intent to set an alarm.
 * Requires Calendar object to set time and days repeating.
 */
fun setAlarm(
    context: Context, date: Calendar,
    message: String = "Sunrise Alarm",
    skipUi: Boolean = false
) {
    val intent = Intent(AlarmClock.ACTION_SET_ALARM).apply {
        putExtra(AlarmClock.EXTRA_HOUR, date.get(Calendar.HOUR))
        putExtra(AlarmClock.EXTRA_MINUTES, date.get(Calendar.MINUTE))
        putExtra(AlarmClock.EXTRA_DAYS, arrayListOf(date.get(Calendar.DAY_OF_WEEK)))
        if (message.isNotEmpty()) {
            putExtra(AlarmClock.EXTRA_MESSAGE, message)
        }
        putExtra(AlarmClock.EXTRA_SKIP_UI, skipUi)
    }
    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
    context.startActivity(intent)
}

/**
 * Launches an intent to the alarm screen of the default clock app.
 */
fun manageAlarms() {
    val intent = Intent(AlarmClock.ACTION_SHOW_ALARMS)
    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
    SunriseApp.getAppContext().startActivity(intent)
}

/**
 * Logic for handling setting alarm to future day (week wrapping).
 */
@RequiresApi(Build.VERSION_CODES.O)
fun getCalendarDate(day: String): Calendar {
    val alarmDate = Calendar.getInstance()

    val today = alarmDate.get(Calendar.DAY_OF_WEEK)
    val future = days.indexOf(day) + 1

    val diff = when {
        future == today -> 7
        future > today  -> future - today
        else            -> 7 - today + future
    }

    alarmDate.add(Calendar.DAY_OF_MONTH, diff)

    return alarmDate
}