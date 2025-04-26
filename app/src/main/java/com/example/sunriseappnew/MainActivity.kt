package com.example.sunriseappnew

import android.Manifest
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.core.app.ActivityCompat
import com.example.sunriseappnew.model.LocationService
import com.example.sunriseappnew.ui.theme.SunriseAppNewTheme
import com.example.sunriseappnew.view.BottomCenteredButton
import com.example.sunriseappnew.view.CheckboxScreen
import com.example.sunriseappnew.viewmodel.SunriseViewModel
import kotlinx.coroutines.delay
import java.util.Calendar
import java.util.Date

class MainActivity : ComponentActivity() {
    private val viewModel: SunriseViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            SunriseAppNewTheme {
                CheckboxScreen( { day -> viewModel.toggleDay(day) }, viewModel.selectedDays )
                BottomCenteredButton  { viewModel.getLocation() }
            }
            GetLocationWithDelay()
        }

        if (!LocationService.hasLocationPermission()) {
            requestLocationPermissions()
        }
    }

    @Composable
    fun GetLocationWithDelay() {
        LaunchedEffect(Unit) {
            viewModel.getLocation()
            delay(1000)
            viewModel.getLocation()

            var date = Date()
            val c: Calendar = Calendar.getInstance()
            c.time = date
            c.add(Calendar.DATE, 1)
            date= c.time
            c.time = date

            Log.d("sunrise", "${LocationService.getNextSunrise(viewModel.location.value, c).time}")
        }

    }

    private fun requestLocationPermissions() {
        Log.d("location", "requesting")
        ActivityCompat.requestPermissions(
            this,
            arrayOf<String>(
                Manifest.permission.ACCESS_FINE_LOCATION,
                Manifest.permission.ACCESS_COARSE_LOCATION
            ),
            1
        )
    }
}