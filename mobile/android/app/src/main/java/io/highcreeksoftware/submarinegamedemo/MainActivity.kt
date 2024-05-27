package io.highcreeksoftware.submarinegamedemo

import android.app.Activity
import android.content.Context
import android.os.Bundle
import android.util.AttributeSet
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.Layout
import androidx.compose.ui.tooling.preview.Preview
import go.Seq
import io.highcreeksoftware.submarinegame.mobile.EbitenView
import io.highcreeksoftware.submarinegamedemo.ui.theme.SubmarineGameDemoTheme

class MainActivity : AppCompatActivity() {
//    override fun onCreate(savedInstanceState: Bundle?) {
//        super.onCreate(savedInstanceState)
//        enableEdgeToEdge()
//        setContent {
//            SubmarineGameDemoTheme {
//                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
//                    Greeting(
//                        name = "Android",
//                        modifier = Modifier.padding(innerPadding)
//                    )
//                }
//            }
//        }
//    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Log.i("MainAct", "Setting content view")
        setContentView(R.layout.main)
        Seq.setContext(applicationContext)
        Log.i("MainAct", "On create is complete")
    }

//    override fun onCreateView(
//        parent: View?,
//        name: String,
//        context: Context,
//        attrs: AttributeSet
//    ): View? {
//        val view = LayoutInflater.from(context).inflate(R.layout.main, parent as ViewGroup)
//
//        return view
//    }

    override fun onPause() {
        super.onPause()
        Log.i("MainAct", "On Pause called")
        findViewById<EbitenView>(R.id.ebitenview).suspendGame()
    }

    override fun onResume() {
        super.onResume()
        Log.i("MainAct", "On Resume called")
        findViewById<EbitenView>(R.id.ebitenview).resumeGame()
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello $name!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    SubmarineGameDemoTheme {
        Greeting("Android")
    }
}