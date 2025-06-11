package file

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>pgdump-mapper</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="bg-light">
    <div class="ms-5 me-5 my-5">
        <div class="row mb-4">
            <div class="col-md-4">
                <h5><label for="tableSelector" class="form-label">Select Table</label></h5>
                <select id="tableSelector" class="form-select">
                    {{range $index, $table := .}}
                        <option value="table-{{$index}}">{{$table.name}}</option>
                    {{end}}
                </select>
            </div>
            <div class="col-md-2 d-flex align-items-end">
                <button class="btn btn-primary w-50" onclick="showTable()">Go!</button>
            </div>
        </div>

        {{range $index, $table := .}}
            <div id="table-{{$index}}" class="table-data" style="display: none;">
                <h4>{{$table.schema}}.{{$table.name}}</h4>
                <table class="table table-bordered table-striped">
                    <thead class="table-dark">
                        <tr>
                            {{if $table.columns}}
                                {{range $_, $c := $table.columns}}
                                    <th>{{$c}}</th>
                                {{end}}
                            {{end}}
                        </tr>
                    </thead>
                    <tbody>
                        {{range $_, $d := $table.data}}
                            <tr>
                                {{range $_, $c := $table.columns}}
                                    <td>{{index $d $c}}</td>
                                {{end}}
                            </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        {{end}}
    </div>
</body>
<script>
    function showTable() {
        var selected = document.getElementById("tableSelector").value;
        var tables = document.getElementsByClassName("table-data");

        for (var i = 0; i < tables.length; i++) {
            tables[i].style.display = "none";
        }

        document.getElementById(selected).style.display = "block";
    }
    document.getElementById("table-0").style.display = "block";
</script>
</html>
`
