package io.highcreeksoftware.submarinegamedemo

import android.content.Context
import android.util.AttributeSet
import android.util.Log
import io.highcreeksoftware.submarinegame.mobile.EbitenView
import java.lang.Exception


class EbitenViewWithErrorHandling @JvmOverloads constructor(context: Context, attrs: AttributeSet? = null, defStyle: Int = 0): EbitenView(context, attrs) {

    override fun onErrorOnGameUpdate(e: Exception?) {
        super.onErrorOnGameUpdate(e)

        e?.let {
            Log.e("Ebiten", "error in game", it)
        }
    }
}