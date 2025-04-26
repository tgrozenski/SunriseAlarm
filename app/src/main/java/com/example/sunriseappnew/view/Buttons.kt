package com.example.sunriseappnew.view

import android.util.Log
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun BottomCenteredButton(onSyncClicked: () -> Unit) {
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.BottomCenter
    ) {
        CustomButton(
            onSyncClicked = onSyncClicked,
            modifier = Modifier.padding(bottom = 32.dp),
        )
    }
}

@Composable
fun CustomButton(
    onSyncClicked: () -> Unit,
    modifier: Modifier = Modifier,
    text: String = "Sync Location",
    onClick: () -> Unit = { onSyncClicked() }
) {
    Button(
        modifier = modifier,
        onClick = onClick,
    ) {
        Text(text)
    }
}
