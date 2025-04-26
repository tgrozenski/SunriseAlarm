package com.example.sunriseappnew.model

import android.app.Application
import android.content.Context

class SunriseApp : Application() {
    companion object {
        private var instance: SunriseApp? = null

        fun getAppContext(): Context { // Provides application context statically
            return instance?.applicationContext ?:
            throw IllegalStateException("Application not initialized!")
        }
    }

    override fun onCreate() {
        super.onCreate()
        instance = this
    }
}