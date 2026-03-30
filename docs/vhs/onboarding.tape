Output docs/assets/onboarding_demo.gif
Set FontSize 34
Set Width 1200
Set Height 900
Set Padding 40
Set Theme "TokyoNight"
Set TypingSpeed 100ms

# Start the recording
Sleep 1s
Type "cue onboarding"
Enter

# Wait for the welcome message
Sleep 4s

# The prompt is: "Press [Enter] or type 'y' to continue... "
# On Ubuntu, the application is likely dying or buffering. 
# We explicitly type 'y' slowly.
Type "y"
Sleep 2s
Enter

# Wait for Step 1
Sleep 4s

# Prompt: "Would you like to install a stack right now? (You can say no) [y/N] "
Type "y"
Sleep 2s
Enter

# TUI Picker interaction
Sleep 3s
Down@2s 2
Up@2s 1
Sleep 1s
Type "q"

# Final readout
Sleep 10s
