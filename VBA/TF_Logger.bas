Attribute VB_Name = "TF_Logger"
Option Explicit

' ========================================
' TF_Logger Module
' Comprehensive logging system for debugging
' ========================================

Private LogFilePath As String
Private LoggingEnabled As Boolean
Private LoggerInitialized As Boolean

Sub InitializeLogger()
    ' Initialize the logger with a file path
    On Error Resume Next

    LogFilePath = ThisWorkbook.Path & "\TradingSystem_Debug.log"
    LoggingEnabled = True
    LoggerInitialized = True

    ' Write header
    Call WriteLog("========================================")
    Call WriteLog("Trading System Debug Log")
    Call WriteLog("Session started: " & Now)
    Call WriteLog("Excel Version: " & Application.Version)
    Call WriteLog("Workbook: " & ThisWorkbook.Name)
    Call WriteLog("========================================")

    On Error GoTo 0
End Sub

Sub WriteLog(message As String)
    ' Write a message to the log file
    On Error Resume Next

    If Not LoggerInitialized Then
        Call InitializeLogger
    End If

    If Not LoggingEnabled Then Exit Sub

    Dim fileNum As Integer
    Dim timestamp As String

    timestamp = Format(Now, "yyyy-mm-dd hh:nn:ss")

    fileNum = FreeFile
    Open LogFilePath For Append As #fileNum
    Print #fileNum, timestamp & " | " & message
    Close #fileNum

    ' Also write to Immediate window for VBA debugging
    Debug.Print timestamp & " | " & message

    On Error GoTo 0
End Sub

Sub WriteLogError(functionName As String, errorNum As Long, errorDesc As String)
    ' Write an error to the log
    Call WriteLog("ERROR in " & functionName & ": [" & errorNum & "] " & errorDesc)
End Sub

Sub WriteLogSection(sectionName As String)
    ' Write a section header
    Call WriteLog("")
    Call WriteLog("--- " & sectionName & " ---")
End Sub

Function GetLogPath() As String
    ' Return the current log file path
    GetLogPath = LogFilePath
End Function

Sub OpenLogFile()
    ' Open the log file in Notepad
    Dim logPath As String
    logPath = GetLogPath()

    If Dir(logPath) <> "" Then
        Shell "notepad.exe """ & logPath & """", vbNormalFocus
    Else
        MsgBox "Log file not found: " & logPath, vbExclamation
    End If
End Sub

Sub ClearLog()
    ' Clear the log file
    On Error Resume Next
    Kill LogFilePath
    On Error GoTo 0
    Call InitializeLogger
    MsgBox "Log file cleared and reinitialized", vbInformation
End Sub

Sub EnableLogging()
    LoggingEnabled = True
    Call WriteLog("Logging enabled")
End Sub

Sub DisableLogging()
    Call WriteLog("Logging disabled")
    LoggingEnabled = False
End Sub

Function IsLoggingEnabled() As Boolean
    IsLoggingEnabled = LoggingEnabled
End Function
