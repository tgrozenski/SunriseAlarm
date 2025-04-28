package com.example.sunriseappnew.view

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

/**
 * Button components for the app.
 */
@Composable
fun BottomCenteredButton(
    onClick: () -> Unit,
    label: String
) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.BottomCenter
    ) {
        CustomButton(
            text = label,
            onClick = onClick,
            modifier = Modifier.padding(bottom = 50.dp),
        )
    }
}

@Composable
fun CustomButton(
    modifier: Modifier = Modifier,
    text: String = "Sync Location",
    onClick: () -> Unit,
) {
    Button(
        modifier = modifier,
        onClick = onClick,
    ) {
        Text(text)
    }
}
