package com.example.sunriseappnew.view

import androidx.compose.foundation.layout.Column
import androidx.compose.material3.Slider
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier

@Composable
fun TimePicker(timeOffset: Float, onValueChange: (Float) -> Unit) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Slider(
            value = timeOffset,
            onValueChange = { onValueChange(it) },
            valueRange = -30f..30f,
            steps = 60
        )
        Text(text = timeOffset.toInt().toString())
    }
}
