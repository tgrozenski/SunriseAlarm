package com.example.sunriseappnew

import android.Manifest
import android.content.Intent
import android.location.Location
import android.os.Build
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import androidx.annotation.RequiresApi
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.graphics.Color
import androidx.core.app.ActivityCompat
import com.example.sunriseappnew.model.LocationService
import com.example.sunriseappnew.ui.theme.SunriseAppNewTheme
import com.example.sunriseappnew.view.BottomCenteredButton
import com.example.sunriseappnew.view.CheckboxScreen
import com.example.sunriseappnew.viewmodel.SunriseViewModel
import kotlinx.coroutines.delay
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.ui.Alignment
import androidx.compose.runtime.State
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.lifecycle.repeatOnLifecycle
import com.example.sunriseappnew.model.manageAlarms
import com.example.sunriseappnew.view.CustomButton
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import android.net.Uri
import android.provider.Settings
import androidx.compose.material3.MaterialTheme
import androidx.compose.ui.text.font.FontWeight

class MainActivity : ComponentActivity() {
    /**
     * View model instance to be used for location and checkbox.
     */
    private val viewModel: SunriseViewModel by viewModels()

    @RequiresApi(Build.VERSION_CODES.O)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            SunriseAppNewTheme {
                Column(
                    modifier = Modifier.fillMaxSize(),
                    verticalArrangement = Arrangement.SpaceEvenly,
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        text = "Sunrise Alarm",
                        style = MaterialTheme.typography.headlineMedium,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                    Text("Check a day to create a sunrise Alarm.")
                    CheckboxScreen({ day -> viewModel.toggleDay(day) }, viewModel.selectedDays)
                    CustomButton(
                        onClick = { manageAlarms() },
                        text = "Manage Alarms"
                    )
                    LocationStatusBox(viewModel.location)
                }
                BottomCenteredButton (
                    onClick = { openLocationPermissionSettings() },
                    label = "Manage Location Permissions"
                )
            }
            if (!LocationService().hasLocationPermission())
                requestLocationPermissions()
            LocationUpdater()
        }
    }

    /**
     * Displays state of location in [SunriseViewModel]
     */
    @Composable
    fun LocationStatusBox(
        location: State<Location?>,
        modifier: Modifier = Modifier
    ) {
        Box(
            modifier = modifier
                .background(
                    color = when (location.value) {
                        null -> Color.Red
                        else ->  Color.Green
                    }
                )
                .padding(16.dp)
        ) {
            Text(
                text = when (location.value) {
                    null -> "No Location"
                    else ->  "Location Found"
                },
                color = Color.White
            )
        }
    }

    /**
     * Periodically check for location
     */
    @Composable
    fun LocationUpdater(
        intervalMs: Long = 3000L,
    ) {
        val lifecycleOwner = LocalLifecycleOwner.current
        var isActive by remember { mutableStateOf(true) }

        LaunchedEffect(key1 = lifecycleOwner.lifecycle) {
            lifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                while (isActive) {
                    viewModel.getLocation()
                    if (viewModel.location.value != null) {
                        isActive = false
                    } else {
                        delay(intervalMs)
                    }
                }
            }
        }
    }

    /**
     * Launch permission settings for user to adjust manually.
     */
    fun openLocationPermissionSettings() {
        val intent = Intent(Settings.ACTION_APPLICATION_DETAILS_SETTINGS).apply {
            data = Uri.fromParts("package", "com.example.sunriseappnew", null)
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
                putExtra("android.provider.extra.APP_PACKAGE", "com.example.sunriseappnew")
            }
            addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        }
        this.startActivity(intent)
    }

    /**
     * Launch permission request upon startup of the app.
     */
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