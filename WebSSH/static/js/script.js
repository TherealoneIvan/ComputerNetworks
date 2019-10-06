$(document).ready(function () {
    $('form').submit(function(){
        $('#execcmd').click();
        return false;
    });
    $('#execcmd').click(function () {
        let cmd = $('#cmd').val();
        if (cmd.replace(/^\s+/, '').replace(/\s+$/, '') == '') {
            return;
        }
        $('#cmd').val("");
        $.ajax({
            type: 'GET',
            url: '/execute',
            data: {
                "command": cmd
            },
            success: function (data) {
                $('<li>', { class: 'list-group-item bg-dark command', text: cmd }).appendTo('#console');
                let output = $('<li class="list-group-item bg-dark output"></li>');
                let outputList = $('<ul class="list-group"></ul>');
                let i = 0;
                for (i = 0; i < data.list.length; i++) {
                    outputList.append($('<li>', { class: 'list-group-item bg-dark output', text: data.list[i] }))
                }
                output.append(outputList);
                output.appendTo('#console');
            },
            failure: function (data) {
                console.log(data);
            }
        });
    });
});