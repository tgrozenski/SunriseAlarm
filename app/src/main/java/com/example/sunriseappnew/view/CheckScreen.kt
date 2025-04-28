package com.example.sunriseappnew.view

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Checkbox
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp

val days = listOf("Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat")

@Composable
fun CheckboxScreen(
    onCheck: (String) -> Unit,
    selectedDays: List<String>,
) {
    Column {
        DayOfWeekCheckboxRow(
//            modifier = Modifier.padding(top = 50.dp),
            selectedDays = selectedDays,
            onCheck = onCheck,
        )
    }
}

@Composable
fun DayOfWeekCheckboxRow(
    modifier: Modifier = Modifier,
    selectedDays: List<String>,
    onCheck: (String) -> Unit,
) {
    Row(
        modifier = modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.SpaceEvenly,
        verticalAlignment = Alignment.CenterVertically
    ) {
        days.forEach { day ->
            DayCheckbox(
                day = day,
                isChecked = selectedDays.contains(day),
                onCheck = onCheck,
            )
        }
    }
}

@Composable
fun DayCheckbox(
    modifier: Modifier = Modifier,
    day: String,
    isChecked: Boolean,
    onCheck: (String) -> Unit,
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Checkbox(
            checked = isChecked,
            onCheckedChange = { checked -> if (checked) onCheck(day) else onCheck(day) }
        )
        Text(
            text = day,
            textAlign = TextAlign.Center,
        )
    }
}