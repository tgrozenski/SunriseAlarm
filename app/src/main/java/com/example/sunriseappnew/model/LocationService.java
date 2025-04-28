package com.example.sunriseappnew.model;
import android.annotation.SuppressLint;
import android.content.Context;
import android.content.pm.PackageManager;
import android.location.Location;
import androidx.core.content.ContextCompat;
import android.Manifest;
import com.google.android.gms.location.FusedLocationProviderClient;
import com.google.android.gms.location.LocationServices;
import com.luckycatlabs.sunrisesunset.SunriseSunsetCalculator;
import java.util.Calendar;

public class LocationService {
    private final Context context;
    private final FusedLocationProviderClient fusedLocationClient;

    public LocationService() {
        this.context = SunriseApp.Companion.getAppContext(); // Use application context to avoid leaks
        this.fusedLocationClient = LocationServices.getFusedLocationProviderClient(context);
    }

    /**
     * Fetches the last location (single shot).
     */
    @SuppressLint("MissingPermission")
    public void getLastLocation(LocationResultCallback callback) {
        if (!hasLocationPermission()) { // checking for permission here
            callback.onPermissionRequired();
            return;
        }

        fusedLocationClient.getLastLocation()
                .addOnSuccessListener(location -> {
                    if (location != null) {
                        callback.onLocationSuccess(location);
                    } else {
                        callback.onLocationUnavailable("Location not available");
                    }
                })
                .addOnFailureListener(e -> {
                    callback.onLocationError(e.getMessage());
                });
    }


    /**
     * Checks for location permission.
     * @return true if granted, else false
     */
    public boolean hasLocationPermission() {
        return ContextCompat.checkSelfPermission(
                context,
                Manifest.permission.ACCESS_FINE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED || ContextCompat.checkSelfPermission(
                context,
                Manifest.permission.ACCESS_COARSE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED;
    }

    /**
     * Interface for all methods to implement in order to call [getLastLocation].
     */
    public interface LocationResultCallback {
        void onLocationSuccess(Location location);
        void onLocationUnavailable(String reason);
        void onLocationError(String error);
        void onPermissionRequired();
    }

    /**
     * Uses luckycatlabs sunrise/sunset calculator to calculate sunrise, gives time two minutes before as margin for error.
     * @return Calendar object for two minutes before requested sunrise.
     */
    public static Calendar getNextSunrise(Location loc, Calendar date) {

        com.luckycatlabs.sunrisesunset.dto.Location location =
                new com.luckycatlabs.sunrisesunset.dto.Location(loc.getLatitude(), loc.getLongitude());

        SunriseSunsetCalculator calculator =
                new SunriseSunsetCalculator(location, date.getTimeZone());

        Calendar sunrise = calculator.getOfficialSunriseCalendarForDate(date);
        sunrise.add(Calendar.MINUTE, -2); // set a bit earlier to account for inaccuracies.

        return sunrise;
    }
}
