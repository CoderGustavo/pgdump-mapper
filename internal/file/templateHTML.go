package file

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>PGDump-Mapper</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <script>
        function showTable() {
            var selected = document.getElementById("tableSelector").value;
            var tables = document.getElementsByClassName("table-data");

            for (var i = 0; i < tables.length; i++) {
                tables[i].style.display = "none";
            }

            document.getElementById("table-" + selected).style.display = "block";
        }
    </script>
</head>
<body class="bg-light">
    <div class="ms-5 me-5 my-5">
        <!-- Select and Button -->
        <div class="row mb-4">
            <div class="col-md-4">
                <h5><label for="tableSelector" class="form-label">Select Table</label></h5>
                <select id="tableSelector" class="form-select">
                    {{range .}}
                        <option value="{{.schema}}.{{.name}}">{{.name}}</option>
                    {{end}}
                </select>
            </div>
            <div class="col-md-2 d-flex align-items-end">
                <button class="btn btn-primary w-50" onclick="showTable()">Go!</button>
            </div>
        </div>

        {{range .}}
            <div id="table-{{.schema}}.{{.name}}" class="table-data" style="display: none;">
                <h4>{{.schema}}.{{.name}}</h4>
                <table class="table table-bordered table-striped">
                    <thead class="table-dark">
                        <tr>
                            {{range $_, $col := .columns}}
                                <th>{{$col}}</th>
                            {{end}}
                        </tr>
                    </thead>
                    <tbody>
                        {{range $_, $rows := .values}}
                            <tr>
                                {{range $_,$col := $rows}}
                                    <td>{{$col}}</td>
                                {{end}}
                            </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        {{end}}
    </div>
</body>
</html>
`
