$path = "HKCU:\Control Panel\Mouse"

Set-ItemProperty $path MouseSpeed 0
Set-ItemProperty $path MouseThreshold1 0
Set-ItemProperty $path MouseThreshold2 0
Set-ItemProperty $path MouseSensitivity 10

Add-Type @"
using System;
using System.Runtime.InteropServices;
public class Native {
 [DllImport("user32.dll")]
 public static extern bool SystemParametersInfo(
   uint uiAction,
   uint uiParam,
   string pvParam,
   uint fWinIni);
}
"@

[Native]::SystemParametersInfo(0x001A,0,$null,0x01 -bor 0x02)
