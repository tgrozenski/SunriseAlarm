package com.example.sunriseappnew.model;
import android.content.Context;
import android.content.pm.PackageManager;
import android.location.Location;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import android.Manifest;
import android.util.Log;
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
     * Fetches the last known location (single shot)
     */
    public void getLastKnownLocation(LocationResultCallback callback) {
        if (!hasLocationPermission()) {
            callback.onPermissionRequired();
            return;
        }

        if (ActivityCompat.checkSelfPermission(
                context,
                android.Manifest.permission.ACCESS_FINE_LOCATION) != PackageManager.PERMISSION_GRANTED
                && ActivityCompat.checkSelfPermission(context,
                android.Manifest.permission.ACCESS_COARSE_LOCATION) != PackageManager.PERMISSION_GRANTED) {
            Log.d("Location", "don't have permission");
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


    public static boolean hasLocationPermission() {
        Context context = SunriseApp.Companion.getAppContext();
        return ContextCompat.checkSelfPermission(
                context,
                Manifest.permission.ACCESS_FINE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED || ContextCompat.checkSelfPermission(
                context,
                Manifest.permission.ACCESS_COARSE_LOCATION
        ) == PackageManager.PERMISSION_GRANTED;
    }

    public interface LocationResultCallback {
        void onLocationSuccess(Location location);
        void onLocationUnavailable(String reason);
        void onLocationError(String error);
        void onPermissionRequired();
    }

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
