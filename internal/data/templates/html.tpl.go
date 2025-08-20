package templates

const HTML = `<!DOCTYPE html>
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
                        <option value="table-{{$index}}">{{$table.Name}}</option>
                    {{end}}
                </select>
            </div>
            <div class="col-md-2 d-flex align-items-end">
                <button class="btn btn-primary w-50" onclick="showTable()">Go!</button>
            </div>
        </div>

        {{range $index, $table := .}}
            <div id="table-{{$index}}" class="table-data table-responsive" style="display: none;">
                <h4>{{$table.Schema}}.{{$table.Name}}</h4>
                <table class="table w-auto table-bordered table-striped">
                    <thead class="table-dark">
                        <tr>
                            {{range $ci, $c := $table.Columns}}
                                <th class="text-nowrap">
                                    {{$c}}
                                    <button class="btn btn-sm btn-outline-secondary btn-filter ms-2" data-table="{{$index}}" data-column="{{$c}}" data-colindex="{{$ci}}" onclick="toggleSelectFilter(this)">
                                        &#128269;
                                    </button>
                                    <select class="form-select form-select-sm filter-select mt-1 d-none" data-table="{{$index}}" data-column="{{$c}}" data-colindex="{{$ci}}">
                                        <option value="">-- All --</option>
                                    </select>
                                </th>
                            {{end}}
                        </tr>
                    </thead>
                    <tbody>
                        {{range $_, $d := $table.Data}}
                            <tr>
                                {{range $_, $c := $table.Columns}}
                                    <td class="text-nowrap">{{index $d $c}}</td>
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

        var table = document.getElementById(selected);
        table.style.display = "block";

        cleanSingleOptionFilters(table);
    }

    function cleanSingleOptionFilters(table) {
        const buttons = table.querySelectorAll(".btn-filter");

        buttons.forEach(button => {
            const tableIndex = button.getAttribute("data-table");
            const colIndex = parseInt(button.getAttribute("data-colindex"));
            const select = button.nextElementSibling;

            const rows = document.querySelectorAll("#table-" + tableIndex + " tbody tr");
            const values = new Set();

            rows.forEach(row => {
                const cell = row.children[colIndex];
                if (cell) values.add(cell.textContent.trim());
            });

            if (values.size <= 1) {
                button.style.display = "none";
                select.style.display = "none";
            } else {
                button.style.display = "inline-block";
                select.classList.add("d-none");

                if (select.options.length <= 1) {
                    Array.from(values).sort().forEach(val => {
                        const option = document.createElement("option");
                        option.value = val;
                        option.textContent = val;
                        select.appendChild(option);
                    });
                }
            }
        });
    }


    function toggleSelectFilter(button) {
        const tableIndex = button.getAttribute("data-table");
        const colIndex = parseInt(button.getAttribute("data-colindex"));
        const select = button.nextElementSibling;

        const rows = document.querySelectorAll("#table-" + tableIndex + " tbody tr");
        const values = new Set();

        rows.forEach(row => {
            const cell = row.children[colIndex];
            if (cell) values.add(cell.textContent.trim());
        });

        if (values.size <= 1) {
            button.style.display = "none";
            select.style.display = "none";
            return;
        }

        button.style.display = "inline-block";

        select.classList.toggle("d-none");

        if (select.options.length > 1) return;

        Array.from(values).sort().forEach(val => {
            const option = document.createElement("option");
            option.value = val;
            option.textContent = val;
            select.appendChild(option);
        });
    }


    document.addEventListener("change", function(event) {
        if (event.target.classList.contains("filter-select")) {
            const tableIndex = event.target.getAttribute("data-table");
            const rows = document.querySelectorAll("#table-" + tableIndex + " tbody tr");

            rows.forEach(row => {
                let visible = true;
                document.querySelectorAll('.filter-select[data-table="' + tableIndex + '"]').forEach(select => {
                    const colIndex = parseInt(select.getAttribute("data-colindex"));
                    const value = select.value.trim();

                    if (value && row.children[colIndex].textContent.trim() !== value) {
                        visible = false;
                    }
                });

                row.style.display = visible ? "" : "none";
            });
        }
    });

    document.getElementById("table-0").style.display = "block";
</script>
</html>
`
