 No extra build step is required, but Windows must see the SDK DLLs when you run. Do this:

  - Open a Developer/MinGW shell in the repo root and put the Orbbec DLL folder on PATH:
      - PowerShell: set PATH=$PWD\orbbec-sdk\OrbbecSDK_v1.10.27\SDK\lib;$PATH
      - Cmd: set PATH=%CD%\orbbec-sdk\OrbbecSDK_v1.10.27\SDK\lib;%PATH%
  - Ensure CGO is on and you’re on Windows: set CGO_ENABLED=1.
  - Drivers: install the Orbbec USB driver from orbbec-sdk/OrbbecSDK_v1.10.27/Driver if not already.
  - Then you can run immediately from the repo root:
      - go run ./cmd/camera-test
      - or tests: go test ./... (unit) and go test -tags=integration ./internal/camera/... (with camera attached).

  You do not need to copy DLLs if PATH is set as above; alternatively, you could copy the four DLLs (OrbbecSDK.dll, depthengine_2_0.dll,
  live555.dll, ob_usb.dll) next to any built exe, but for go run setting PATH is the safest.