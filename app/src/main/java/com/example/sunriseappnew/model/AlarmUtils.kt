package com.example.sunriseappnew.model

import android.content.Context
import android.content.Intent
import android.os.Build
import android.provider.AlarmClock
import android.util.Log
import androidx.annotation.RequiresApi
import java.util.Calendar


val days = listOf("Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat")

fun setAlarm(context: Context, date: Calendar, message: String = "", skipUi: Boolean = false) {
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

fun manageAlarms() {
    val intent = Intent(AlarmClock.ACTION_DISMISS_ALARM).apply {
        putExtra(AlarmClock.EXTRA_ALARM_SEARCH_MODE, AlarmClock.ALARM_SEARCH_MODE_LABEL)
        putExtra(AlarmClock.EXTRA_MESSAGE, "Sunrise Alarm")
    }
    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
    SunriseApp.getAppContext().startActivity(intent)
}

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

    Log.d("sunrise", "diff: ${diff}")
    Log.d("sunrise", "future, today: ${future}, $today")
    Log.d("sunrise", "new alarm: ${alarmDate.get(Calendar.DAY_OF_MONTH)}")
    Log.d("sunrise", "new alarm: ${alarmDate.get(Calendar.DAY_OF_WEEK)}")

    return alarmDate
}