Attribute VB_Name = "PQ_Setup"
Option Explicit

Private Const cNamedRepoDir As String = "RepoDataDir"
Private Const cQueryName As String = "PQ_Screened"
Private Const cSheetName As String = "Screened"
Private Const cTableName As String = "tblScreened"

' Entry point: run once after cloning the repo or if the query breaks
Public Sub InstallOrRepairPowerQuery()
    Dim repoDataDir As String
    repoDataDir = ThisWorkbook.Path & "\data"

    EnsureSettingsCell repoDataDir
    AddOrUpdateQuery
    BindQueryToSheet
    MsgBox "Power Query installed/refreshed.", vbInformation
End Sub

' Ensure we have a cell and a defined name with the data directory
Private Sub EnsureSettingsCell(ByVal repoDataDir As String)
    Dim ws As Worksheet
    On Error Resume Next
    Set ws = ThisWorkbook.Worksheets("Settings")
    On Error GoTo 0
    If ws Is Nothing Then
        Set ws = ThisWorkbook.Worksheets.Add(After:=Sheets(Sheets.Count))
        ws.Name = "Settings"
        ws.Range("A1").Value = "Data directory"
    End If

    ws.Range("B1").Value = repoDataDir
    On Error Resume Next
    ThisWorkbook.Names(cNamedRepoDir).Delete
    On Error GoTo 0
    ThisWorkbook.Names.Add Name:=cNamedRepoDir, RefersTo:=ws.Range("B1")
End Sub

' Add or update the M query that reads data\screened.csv
Private Sub AddOrUpdateQuery()
    Dim m As String
    m = _
"let" & vbCrLf & _
"  DataDir = Excel.CurrentWorkbook(){[Name=""" & cNamedRepoDir & """]}[Content]{0}[Column1]," & vbCrLf & _
"  FilePath = DataDir & ""\screened.csv""," & vbCrLf & _
"  Source = Csv.Document(File.Contents(FilePath),[Delimiter="","", Columns=3, Encoding=65001, QuoteStyle=QuoteStyle.Csv])," & vbCrLf & _
"  Promote = Table.PromoteHeaders(Source, [PromoteAllScalars=true])," & vbCrLf & _
"  Types = Table.TransformColumnTypes(Promote,{{""Ticker"", type text},{""Preset"", type text},{""AsOf"", type date}})," & vbCrLf & _
"  Upper = Table.TransformColumns(Types, {{""Ticker"", Text.Upper, type text}})" & vbCrLf & _
"in" & vbCrLf & _
"  Upper"

    Dim q As WorkbookQuery
    On Error Resume Next
    Set q = ThisWorkbook.Queries(cQueryName)
    On Error GoTo 0

    If q Is Nothing Then
        ThisWorkbook.Queries.Add Name:=cQueryName, Formula:=m
    Else
        q.Formula = m
    End If
End Sub

' Load the query to a sheet table (creates sheet/table if missing)
Private Sub BindQueryToSheet()
    Dim ws As Worksheet, lo As ListObject, conn As WorkbookConnection, cn As OLEDBConnection
    On Error Resume Next
    Set ws = ThisWorkbook.Worksheets(cSheetName)
    On Error GoTo 0
    If ws Is Nothing Then
        Set ws = ThisWorkbook.Worksheets.Add(After:=Sheets(Sheets.Count))
        ws.Name = cSheetName
    End If

    ' Remove previous table if present
    On Error Resume Next
    Set lo = ws.ListObjects(cTableName)
    If Not lo Is Nothing Then lo.Unlist
    On Error GoTo 0
    ws.Cells.Clear

    ' Create a PQ connection-backed table at A1
    Dim src As String
    src = "OLEDB;Provider=Microsoft.Mashup.OleDb.1;Data Source=$Workbook$;Location=" & cQueryName & ";Extended Properties="""""

    Dim qt As QueryTable
    Set qt = ws.ListObjects.Add(SourceType:=0, Source:=src, Destination:=ws.Range("A1")).QueryTable
    qt.CommandType = xlCmdSql
    qt.CommandText = Array("SELECT * FROM [" & cQueryName & "]")

    ws.ListObjects(1).Name = cTableName

    ' Refresh settings
    Set conn = ThisWorkbook.Connections("Query - " & cQueryName)
    Set cn = conn.OLEDBConnection
    cn.BackgroundQuery = False
    cn.RefreshOnFileOpen = True

    ' First refresh
    conn.Refresh
    ws.Columns.AutoFit
End Sub
