package com.example.sunriseappnew.viewmodel
import android.location.Location
import android.os.Build
import android.util.Log
import android.widget.Toast
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
import java.util.Calendar


/**
 * The only viewModel for app, stores the checkbox and location state & data.
 */
class SunriseViewModel : ViewModel() {

    /** Mutable private instance of selected days checkbox.*/
    private val _selectedDays = mutableStateListOf<String>()

    /** Immutable public instance of selected days checkbox.*/
    val selectedDays: List<String> = _selectedDays

    /** Mutable private instance of location. .*/
    private val _location = mutableStateOf<Location?>(null)
    /** Immutable public instance of location. .*/
    val location: State<Location?> get() = _location

    /** Mutable private instance of timeOffset */
    private val _timeOffset = mutableStateOf(0f)
    val timeOffset: State<Float> get() = _timeOffset

    fun setTimeOffset(offset: Float) {
        _timeOffset.value = offset
    }

    /**
     * Uses [LocationService] to get the last known location if location is null and permission
     * is granted, implements LocationServiceCallback.
     */
    fun getLocation() {
        if (_location.value == null && LocationService().hasLocationPermission())
            LocationService().getLastLocation(
               object : LocationResultCallback {
                   override fun onLocationSuccess(loc: Location?) {
                       _location.value = loc
                       Toast.makeText(SunriseApp.getAppContext(), "Location Received", Toast.LENGTH_SHORT).show()
                       Log.d("location", "success: ${_location.value}")
                   }
                   override fun onLocationUnavailable(reason: String?) {
                       Toast.makeText(SunriseApp.getAppContext(),
                           "Location Unavailable, Try opening another app with location services",
                           Toast.LENGTH_LONG).show()
                       Log.d("location", "failure: $reason")
                   }
                   override fun onLocationError(error: String?) {
                       Toast.makeText(SunriseApp.getAppContext(),
                           "Error getting location",
                           Toast.LENGTH_SHORT).show()
                   }
                   override fun onPermissionRequired() {
                       Log.d("location", "permission is missing")
                       Toast.makeText(SunriseApp.getAppContext(),
                           "No permissions, please grant location permissions",
                           Toast.LENGTH_SHORT).show()
                   }
               }
            )
        else
            Log.d("location", "already have location: ${_location.value}")
    }

    /**
     * Toggle day Logic, tries to launch an alarm intent if checked.
     * Updates selectedDays list for UI to observe.
     */
    @RequiresApi(Build.VERSION_CODES.O)
    fun toggleDay(day: String) {
        if (day in selectedDays) {
            _selectedDays.remove(day)
        } else {
            _selectedDays.add(day)

            if (_location.value != null) {
                val sunrise: Calendar = LocationService.getNextSunrise(_location.value, getCalendarDate(day))
                // Adjust to offset here
                sunrise.add(Calendar.MINUTE, timeOffset.value.toInt())
                setAlarm(SunriseApp.getAppContext(), sunrise, skipUi = true)
            }
        }
    }
}