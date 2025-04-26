package com.example.sunriseappnew.viewmodel
import android.location.Location
import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import com.example.sunriseappnew.model.LocationService
import com.example.sunriseappnew.model.LocationService.LocationResultCallback
import com.example.sunriseappnew.model.SunriseApp
import com.example.sunriseappnew.model.setAlarm
import androidx.compose.runtime.State

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

        Log.d("location", "${_location.value}")
        Log.d("location", "${location.value}")

    }

    fun toggleDay(day: String) {
        if (day in selectedDays) {
            _selectedDays.remove(day)
        } else {
            _selectedDays.add(day)
            setAlarm(SunriseApp.getAppContext(), 8, 20, skipUi = true)
        }
        Log.d("checks", "$selectedDays")
    }
}
