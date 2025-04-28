package com.example.sunriseappnew.viewmodel
import android.location.Location
import android.os.Build
import android.util.Log
import androidx.annotation.RequiresApi
import androidx.compose.runtime.State
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import com.example.sunriseappnew.model.LocationService
import com.example.sunriseappnew.model.LocationService.LocationResultCallback
import com.example.sunriseappnew.model.SunriseApp
import com.example.sunriseappnew.model.getCalendarDate
import com.example.sunriseappnew.model.setAlarm


class SunriseViewModel : ViewModel() {

    private val _selectedDays = mutableStateListOf<String>()
    val selectedDays: List<String> = _selectedDays

    private val _location = mutableStateOf<Location?>(null)
    val location: State<Location?> get() = _location

    private val locationService = LocationService()

    private val locationImplementation = object : LocationResultCallback {
        override fun onLocationSuccess(loc: Location?) {
            _location.value = loc
            Log.d("location", "success: ${_location.value}")
            Log.d("location", "success: ${location.value}")
        }
        override fun onLocationUnavailable(reason: String?) {
            Log.d("location", "failure: $reason")
        }
        override fun onLocationError(error: String?) {
            Log.d("location", "error: $error")
        }
        override fun onPermissionRequired() {
            Log.d("location", "permission is missing")
        }
    }

    fun getLocation() {
        if (_location.value == null)
            locationService.getLastKnownLocation(locationImplementation)
        else
            Log.d("location", "already have location: ${_location.value}")
    }

    @RequiresApi(Build.VERSION_CODES.O)
    fun toggleDay(day: String) {
        if (day in selectedDays) {
            _selectedDays.remove(day)
        } else {
            _selectedDays.add(day)

            if (_location.value != null) {
                val sunrise = LocationService.getNextSunrise(_location.value, getCalendarDate(day))
                Log.d("sunrise", "sunrise time: ${sunrise.time}")
                setAlarm(SunriseApp.getAppContext(), sunrise, message = "Sunrise Alarm", skipUi = true)
            } else {
                Log.d("location", "no location to use")
            }
        }
    }
}