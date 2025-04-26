package com.example.sunriseappnew.model

import android.content.Context
import android.content.Intent
import android.provider.AlarmClock
import java.util.Calendar


fun setAlarm(context: Context, hour: Int, minute: Int, message: String = "", skipUi: Boolean = false) {
    val intent = Intent(AlarmClock.ACTION_SET_ALARM).apply {
        putExtra(AlarmClock.EXTRA_HOUR, hour)
        putExtra(AlarmClock.EXTRA_MINUTES, minute)
        putExtra(AlarmClock.EXTRA_DAYS, arrayListOf(Calendar.SUNDAY))
        if (message.isNotEmpty()) {
            putExtra(AlarmClock.EXTRA_MESSAGE, message)
        }
        putExtra(AlarmClock.EXTRA_SKIP_UI, skipUi)
    }
    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
    context.startActivity(intent)
}
