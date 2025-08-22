package com.example.sunriseappnew.view

import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowDown
import androidx.compose.material.icons.filled.KeyboardArrowUp
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalSoftwareKeyboardController
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import kotlin.math.abs

@Composable
fun SignedNumberInput(
    value: Int,
    onValueChange: (Int) -> Unit,
    modifier: Modifier = Modifier
) {
    var text by remember(value) { mutableStateOf(abs(value).toString()) }
    val keyboardController = LocalSoftwareKeyboardController.current

    Row(
        verticalAlignment = Alignment.CenterVertically,
        modifier = modifier
    ) {
        IconButton(onClick = {
            if (value != 0) {
                onValueChange(value * -1)
            }
        }) {
            Icon(
                imageVector = if (value >= 0) Icons.Default.KeyboardArrowDown else Icons.Default.KeyboardArrowUp,
                contentDescription = "Toggle Sign"
            )
        }

        OutlinedTextField(
            value = text,
            onValueChange = { newText ->
                val filteredText = newText.filter { it.isDigit() }
                val newValue = filteredText.toIntOrNull() ?: 0

                // Enforce the 0-30 limit for the input field
                if (newValue <= 60) {
                    text = filteredText
                    val signedValue = if (value >= 0) newValue else -newValue
                    onValueChange(signedValue)
                }
            },
            label = { Text("Offset") },
            keyboardOptions = KeyboardOptions(
                keyboardType = KeyboardType.Number,
                imeAction = ImeAction.Done
            ),
            keyboardActions = KeyboardActions(
                onDone = {
                    keyboardController?.hide()
                }
            ),
            modifier = Modifier.weight(1f).padding(end = 8.dp)
        )
    }
}
