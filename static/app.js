/*
 * Secondly: Configuration manager for Go language apps
 * Copyright (c) 2015 Gregory Eremin
 *
 * Source: https://github.com/localhots/secondly
 * Licence: https://github.com/localhots/secondly/blob/master/LICENSE
 */

function loadFields(callback) {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/fields.json", true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var fields = JSON.parse(xhr.responseText);
                callback(fields);
            }
        }
    };
    xhr.send(null);
}

function saveFields(payload, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/save", true);
    xhr.setRequestHeader("Content-Type", "application/json; charset=utf-8");
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                callback(JSON.parse(xhr.responseText));
            } else {
                callback({"success": false, "msg": "Failed to save config"})
            }
        }
    };
    xhr.send(JSON.stringify(payload));
}

function drawForm(fields) {
    var elems = [];

    var curLevel = 0;
    var titlesPrinted = {};
    for (var i = 0; i < fields.length; i++) {
        var field = fields[i];

        var tokens = field.path.split(".");
        var section = tokens.slice(0, -1).join(".");

        if (section != "" && !titlesPrinted[section]) {
            titlesPrinted[section] = 1;
            elems.push({
                level: tokens.length - 1,
                nodes: [makeSectionNode("/"+ section)],
            });
        }

        elems.push({
            level: tokens.length - 1,
            nodes: makeFieldNode(field),
        });
    }

    render(elems);
}

function render(elems) {
    var fields = document.getElementById("fields");

    for (var i = 0; i < elems.length; i++) {
        var row = elems[i];
        var nodes = row.nodes;

        for (var j = 0; j < row.level; j++) {
            nodes.unshift(makePaddingNode())
        }

        fields.appendChild(makeRow(nodes));
    }
}

function makePaddingNode() {
    var div = document.createElement("div");
    div.setAttribute("class", "padding");
    return div;
}

function makeRow(nodes) {
    var div = document.createElement("div");
    div.setAttribute("class", "row");

    for (var i = 0; i < nodes.length; i++) {
        div.appendChild(nodes[i]);
    }

    return div;
}

function makeFieldNode(field) {
    var formGroup = [],
        label = makeLabelNode(field.path, field.name),
        input = document.createElement("input");


    input.setAttribute("id", field.path);

    if (field.kind !== "bool") {
        input.value = field.value;
    } else {
        if (field.value) {
            input.setAttribute("checked", "checked");
        }
    }

    input.setAttribute("data-type", field.kind);

    switch (field.kind) {
    case "string":
        input.setAttribute("type", "text");
        formGroup.push(label);
        formGroup.push(input);
        break;
    case "bool":
        input.setAttribute("type", "checkbox");
        label.innerHTML = "";
        label.appendChild(document.createTextNode(field.name));
        formGroup.push(label);
        formGroup.push(input);
        break;
    case "int":
    case "int8":
    case "int16":
    case "int32":
    case "int64":
    case "uint":
    case "uint8":
    case "uint16":
    case "uint32":
    case "uint64":
    case "float32":
    case "float64":
        input.setAttribute("type", "number");
        switch (field.kind) {
        case "int8":
            input.setAttribute("min", "-128");
            input.setAttribute("max", "127");
            break;
        case "int16":
            input.setAttribute("min", "-32768");
            input.setAttribute("max", "32767");
            break;
        case "int32":
            input.setAttribute("min", "-2147483648");
            input.setAttribute("max", "2147483647");
            break;
        case "int": // Assuming x86-64 architecture
        case "int64":
            input.setAttribute("min", "-9223372036854775808");
            input.setAttribute("max", "9223372036854775807");
            break;
        case "uint8":
            input.setAttribute("min", "0");
            input.setAttribute("max", "255");
            break;
        case "uint16":
            input.setAttribute("min", "0");
            input.setAttribute("max", "65535");
            break;
        case "uint32":
            input.setAttribute("min", "0");
            input.setAttribute("max", "4294967295");
            break;
        case "uint": // Assuming x86-64 architecture
        case "uint64":
            input.setAttribute("min", "0");
            input.setAttribute("max", "18446744073709551615");
            break;
        case "float32":
        case "float64":
            input.setAttribute("step", "any");
            break;
        }
        formGroup.push(label);
        formGroup.push(input);
        break;
    default:
        console.log("Invalid field type: "+ field.kind, field.path)
    }

    return formGroup;
}

function makeSectionNode(section) {
    var h2 = document.createElement("h2"),
        contents = document.createTextNode(section);
    h2.appendChild(contents);
    return h2;
}

function makeDivNode(classes) {
    var div = document.createElement("div");
    div.setAttribute("class", classes);
    return div;
}

function makeLabelNode(forId, text) {
    var label = document.createElement("label"),
        contents = document.createTextNode(text);
    label.setAttribute("for", forId);
    label.appendChild(contents);
    return label;
}

function makePayload(elems) {
    var payload = {};
    for (path in elems) {
        var value = elems[path],
            tokens = path.split('.'),
            parents = tokens.slice(0, -1),
            key = tokens.slice(-1)[0],
            parent = payload;

        for (var i = 0; i < parents.length; i++) {
            var pkey = parents[i];

            if (!parent[pkey]) {
                parent[pkey] = {}
            }
            parent = parent[pkey];
        }

        parent[key] = value;
    }

    return payload;
}

document.getElementById("config").addEventListener("submit", function(e){
    e.preventDefault();

    var elems = {},
        inputs = document.getElementsByTagName("input");

    for (var i = 0; i < inputs.length; i++) {
        var input = inputs[i],
            type = input.getAttribute("data-type"),
            path = input.getAttribute("id"),
            value = input.value;

        switch (type) {
        case "string":
            elems[path] = value;
            break;
        case "bool":
            elems[path] = input.checked;
            break;
        case "int":
        case "int8":
        case "int16":
        case "int32":
        case "int64":
        case "uint":
        case "uint8":
        case "uint16":
        case "uint32":
        case "uint64":
            elems[path] = parseInt(value, 10);
            break;
        case "float32":
        case "float64":
            elems[path] = parseFloat(value);
            break;
        }
    }

    saveFields(makePayload(elems), function(resp){
        var notice = document.getElementById("notice");
        notice.innerHTML = resp.msg;
        if (resp.success) {
            notice.setAttribute("class", "success");
        } else {
            notice.setAttribute("class", "error");
        }
        notice.style.display = "block";
        window.setTimeout(function() {
            notice.style.display = "none";
        }, 2000);
    });

    return false;
});

loadFields(drawForm);
