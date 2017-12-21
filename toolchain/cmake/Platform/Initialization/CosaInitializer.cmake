write_sep()

# Set up compiler flags
include(SetupCompilerSettings)
include(SetupArduinoSettings)

# Get the `cosa` version
include(CosaDetectVersion)

# Register paths
include(CosaRegisterHardwarePlatform)

# Find examples, libraries, and programs
include(CosaFindPrograms)
