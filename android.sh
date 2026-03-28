#! /bin/bash

ebitenmobile bind -target android -ldflags="max-page-size=16384" -javapkg io.highcreeksoftware.submarinegame -o ./mobile/android/libs/submarinegame.aar ./cmd/mobile/
