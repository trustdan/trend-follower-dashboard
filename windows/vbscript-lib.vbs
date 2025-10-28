' ===========================================================================
' vbscript-lib.vbs - Reusable VBScript Functions for Excel UI Generation
' ===========================================================================
' Purpose: Helper functions for automated worksheet creation
' Used by: 1-setup-all.bat during workbook generation
' Created: M22 - Automated UI Generation
' ===========================================================================

' ---------------------------------------------------------------------------
' Worksheet Creation and Formatting
' ---------------------------------------------------------------------------

' CreateWorksheet - Creates a new worksheet with name and tab color
' Parameters:
'   wb        - Workbook object
'   sheetName - Name for the new sheet
'   tabColor  - RGB color for tab (use RGB(r,g,b))
' Returns: Worksheet object
' ---------------------------------------------------------------------------
Function CreateWorksheet(wb, sheetName, tabColor)
    Dim ws
    Set ws = wb.Worksheets.Add(, wb.Worksheets(wb.Worksheets.Count))
    ws.Name = sheetName
    ws.Tab.Color = tabColor
    Set CreateWorksheet = ws
End Function

' ---------------------------------------------------------------------------
' Cell Formatting
' ---------------------------------------------------------------------------

' AddLabel - Adds formatted text label to a cell
' Parameters:
'   ws    - Worksheet object
'   cell  - Cell reference (e.g., "A1")
'   text  - Text to display
'   bold  - True/False for bold font
'   size  - Font size (points)
'   color - Font color RGB value
' ---------------------------------------------------------------------------
Sub AddLabel(ws, cell, text, bold, size, color)
    ws.Range(cell).Value = text
    ws.Range(cell).Font.Bold = bold
    ws.Range(cell).Font.Size = size
    ws.Range(cell).Font.Color = color
End Sub

' AddFormula - Adds formula to a cell with formatting
' Parameters:
'   ws         - Worksheet object
'   cell       - Cell reference (e.g., "B4")
'   formula    - Excel formula (with = prefix)
'   numFormat  - Number format string (e.g., "$#,##0.00")
' ---------------------------------------------------------------------------
Sub AddFormula(ws, cell, formula, numFormat)
    ws.Range(cell).Formula = formula
    If numFormat <> "" Then
        ws.Range(cell).NumberFormat = numFormat
    End If
End Sub

' SetCellFormat - Sets number format for a cell or range
' Parameters:
'   ws        - Worksheet object
'   cellRange - Cell or range reference (e.g., "B4" or "B4:B10")
'   format    - Number format string
' ---------------------------------------------------------------------------
Sub SetCellFormat(ws, cellRange, format)
    ws.Range(cellRange).NumberFormat = format
End Sub

' ---------------------------------------------------------------------------
' Controls - Buttons, Dropdowns, Checkboxes
' ---------------------------------------------------------------------------

' AddButton - Adds a button control to worksheet
' Parameters:
'   ws        - Worksheet object
'   cell      - Cell reference for positioning
'   caption   - Button text
'   macroName - VBA macro to run (e.g., "TFEngine.CalculatePositionSize")
'   width     - Button width in pixels (optional, default 120)
'   height    - Button height in pixels (optional, default 30)
' ---------------------------------------------------------------------------
Sub AddButton(ws, cell, caption, macroName, width, height)
    Dim btn, w, h
    w = width
    h = height
    If w = 0 Then w = 120
    If h = 0 Then h = 30

    Set btn = ws.Buttons.Add(ws.Range(cell).Left, ws.Range(cell).Top, w, h)
    btn.Text = caption
    btn.OnAction = macroName
    btn.Font.Bold = True
End Sub

' AddDropdown - Adds data validation dropdown to a cell
' Parameters:
'   ws    - Worksheet object
'   cell  - Cell reference (e.g., "B8")
'   items - Comma-separated list of values (e.g., "stock,opt-delta-atr,opt-maxloss")
' ---------------------------------------------------------------------------
Sub AddDropdown(ws, cell, items)
    Dim validation
    On Error Resume Next
    ws.Range(cell).Validation.Delete
    On Error Goto 0

    Set validation = ws.Range(cell).Validation
    validation.Add 3, 1, 1, items ' Type=xlValidateList, AlertStyle=xlValidAlertStop, Operator=xlBetween
    validation.IgnoreBlank = True
    validation.InCellDropdown = True
End Sub

' AddCheckbox - Adds a checkbox control to worksheet
' Parameters:
'   ws       - Worksheet object
'   cell     - Cell reference for positioning
'   caption  - Checkbox label text
'   width    - Width in pixels (optional, default 300)
' Returns: OLEObject (checkbox control)
' ---------------------------------------------------------------------------
Function AddCheckbox(ws, cell, caption, width)
    Dim chk, w
    w = width
    If w = 0 Then w = 300

    Set chk = ws.OLEObjects.Add("Forms.CheckBox.1", False, False, _
        ws.Range(cell).Left, ws.Range(cell).Top, w, 20)
    chk.Object.Caption = caption
    chk.Object.Value = False
    Set AddCheckbox = chk
End Function

' ---------------------------------------------------------------------------
' Area Formatting
' ---------------------------------------------------------------------------

' FormatResultArea - Formats a range as result display area
' Parameters:
'   ws        - Worksheet object
'   startCell - Top-left cell (e.g., "A15")
'   endCell   - Bottom-right cell (e.g., "B22")
'   bgColor   - Background color RGB value (0 for no background)
'   border    - True to add border
' ---------------------------------------------------------------------------
Sub FormatResultArea(ws, startCell, endCell, bgColor, border)
    Dim rng
    Set rng = ws.Range(startCell & ":" & endCell)

    If bgColor <> 0 Then
        rng.Interior.Color = bgColor
    End If

    If border Then
        rng.Borders.LineStyle = 1  ' xlContinuous
        rng.Borders.Weight = 2     ' xlThin
    End If
End Sub

' FormatHeader - Formats a range as header area
' Parameters:
'   ws        - Worksheet object
'   cellRange - Cell or range reference
'   bgColor   - Background color RGB value
'   fontColor - Font color RGB value
' ---------------------------------------------------------------------------
Sub FormatHeader(ws, cellRange, bgColor, fontColor)
    Dim rng
    Set rng = ws.Range(cellRange)
    rng.Interior.Color = bgColor
    rng.Font.Color = fontColor
    rng.Font.Bold = True
    rng.Borders.LineStyle = 1
End Sub

' ---------------------------------------------------------------------------
' Column Width and Row Height
' ---------------------------------------------------------------------------

' SetColumnWidth - Sets width for one or more columns
' Parameters:
'   ws      - Worksheet object
'   columns - Column letter or range (e.g., "A" or "A:C")
'   width   - Width in characters
' ---------------------------------------------------------------------------
Sub SetColumnWidth(ws, columns, width)
    ws.Columns(columns).ColumnWidth = width
End Sub

' SetRowHeight - Sets height for one or more rows
' Parameters:
'   ws     - Worksheet object
'   rows   - Row number or range (e.g., "1" or "1:5")
'   height - Height in points
' ---------------------------------------------------------------------------
Sub SetRowHeight(ws, rows, height)
    ws.Rows(rows).RowHeight = height
End Sub

' ---------------------------------------------------------------------------
' Protection and Lock
' ---------------------------------------------------------------------------

' LockCell - Locks a cell (requires sheet protection to take effect)
' Parameters:
'   ws   - Worksheet object
'   cell - Cell reference
' ---------------------------------------------------------------------------
Sub LockCell(ws, cell)
    ws.Range(cell).Locked = True
End Sub

' UnlockCell - Unlocks a cell to allow editing when sheet is protected
' Parameters:
'   ws   - Worksheet object
'   cell - Cell reference
' ---------------------------------------------------------------------------
Sub UnlockCell(ws, cell)
    ws.Range(cell).Locked = False
End Sub

' ---------------------------------------------------------------------------
' Utility Functions
' ---------------------------------------------------------------------------

' MergeCells - Merges a range of cells
' Parameters:
'   ws        - Worksheet object
'   startCell - Top-left cell
'   endCell   - Bottom-right cell
' ---------------------------------------------------------------------------
Sub MergeCells(ws, startCell, endCell)
    ws.Range(startCell & ":" & endCell).Merge
End Sub

' CenterAlign - Center-aligns cell content
' Parameters:
'   ws        - Worksheet object
'   cellRange - Cell or range reference
' ---------------------------------------------------------------------------
Sub CenterAlign(ws, cellRange)
    ws.Range(cellRange).HorizontalAlignment = -4108  ' xlCenter
End Sub

' RightAlign - Right-aligns cell content
' Parameters:
'   ws        - Worksheet object
'   cellRange - Cell or range reference
' ---------------------------------------------------------------------------
Sub RightAlign(ws, cellRange)
    ws.Range(cellRange).HorizontalAlignment = -4152  ' xlRight
End Sub

' AddBorder - Adds border to a range
' Parameters:
'   ws        - Worksheet object
'   cellRange - Cell or range reference
'   weight    - Border weight: 1=Hairline, 2=Thin, 3=Medium, 4=Thick
' ---------------------------------------------------------------------------
Sub AddBorder(ws, cellRange, weight)
    Dim rng
    Set rng = ws.Range(cellRange)
    rng.Borders.LineStyle = 1  ' xlContinuous
    If weight > 0 Then
        rng.Borders.Weight = weight
    End If
End Sub

' ---------------------------------------------------------------------------
' Named Ranges
' ---------------------------------------------------------------------------

' AddNamedRange - Creates a workbook-level named range
' Parameters:
'   wb        - Workbook object
'   rangeName - Name for the range
'   ws        - Worksheet object
'   cellRange - Cell or range reference
' ---------------------------------------------------------------------------
Sub AddNamedRange(wb, rangeName, ws, cellRange)
    On Error Resume Next
    wb.Names.Add rangeName, ws.Range(cellRange)
    On Error Goto 0
End Sub

' ===========================================================================
' End of vbscript-lib.vbs
' ===========================================================================
