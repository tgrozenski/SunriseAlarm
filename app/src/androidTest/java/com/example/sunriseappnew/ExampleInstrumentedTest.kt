package com.example.sunriseappnew

import android.location.Location
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import com.example.sunriseappnew.model.LocationService
import com.example.sunriseappnew.model.getCalendarDate
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.runner.RunWith
import java.util.Calendar
import kotlin.math.abs


/**
 * Instrumented test, which will execute on an Android device.
 *
 * See [testing documentation](http://d.android.com/tools/testing).
 */
@RunWith(AndroidJUnit4::class)
class ExampleInstrumentedTest {
//    @Test
    fun useAppContext() {
        // Context of the app under test.
        val appContext = InstrumentationRegistry.getInstrumentation().targetContext
        assertEquals("com.example.sunriseappnew", appContext.packageName)
    }

//    @Test
    fun testNextSunrise() {
        val targetLocation = Location("")
        targetLocation.latitude = 34.052235 // Coordinates for LA
        targetLocation.longitude = -118.243683
        val calendar = Calendar.getInstance()
        calendar.set(Calendar.DAY_OF_MONTH, 26)

        println("This is the calendar instance ${calendar.time}")
        val sunrise = LocationService.getNextSunrise(targetLocation, calendar)
        println("This is the sunrise instance ${sunrise.time}")

        assertEquals(6, sunrise.time.hours)
        assertTrue(abs(8 - sunrise.time.minutes) <= 1);
    }

//    @Test
    fun testNextSunrise2() {
        val targetLocation = Location("")
        targetLocation.latitude = 34.052235
        targetLocation.longitude = -118.243683
        val calendar = Calendar.getInstance()
        calendar.set(Calendar.DAY_OF_MONTH, 27)

        println("This is the calendar instance ${calendar.time}")
        val sunrise = LocationService.getNextSunrise(targetLocation, calendar)
        println("This is the sunrise instance ${sunrise.time}")

        assertEquals(6, sunrise.time.hours)
        assertTrue(abs(7 - sunrise.time.minutes) <= 1);
    }

//    @Test
    fun testNextSunrise3() {
        val targetLocation = Location("")
        targetLocation.latitude = 34.052235
        targetLocation.longitude = -118.243683
        val calendar = Calendar.getInstance()

        calendar.set(Calendar.DAY_OF_MONTH, 28)

        println("This is the calendar instance ${calendar.time}")
        val sunrise = LocationService.getNextSunrise(targetLocation, calendar)
        println("This is the sunrise instance ${sunrise.time}")

        assertEquals(6, sunrise.time.hours)
        assertTrue(abs(6 - sunrise.time.minutes) <= 1);
    }

//    @Test
    fun testNextSunrise4() {
        val targetLocation = Location("")
        targetLocation.latitude = 34.052235
        targetLocation.longitude = -118.243683


        val calendar = Calendar.getInstance()
        calendar.set(Calendar.MONTH, 6)
        calendar.set(Calendar.DAY_OF_MONTH, 4)

        println("This is the calendar instance ${calendar.time}")
        val sunrise = LocationService.getNextSunrise(targetLocation, calendar)
        println("This is the sunrise instance ${sunrise.time}")

        assertEquals(5, sunrise.time.hours)
        assertTrue(abs(46 - sunrise.time.minutes) <= 1);
    }


//    @Test
    fun testGetCalendarDate() {
        getCalendarDate("Tue")
    }
}