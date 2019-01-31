// Convert and HTML special char to the UTF-8
function convertHTMLEntity(text){
    const span = document.createElement('span');

    return text
    .replace(/&[#A-Za-z0-9]+;/gi, (entity,position,text)=> {
        span.innerHTML = entity;
        return span.innerText;
    });
}


// Reload note
function reloadNote (lastSave) {
    $.getJSON( "/api/"+path, function( data ) {
        // Ignore load if last save < 2
        if (lastSave < 2) return;

        $("#loader").hide();
        $("#errorMessage").hide();
        $("#noteField").val(data['content']);
        $("#noteField").prop("disabled", false);

        if(data['updateAtHumanize'] === '') {
            $("#updateAtHumanize").parent(".ribbon").hide();
        } else {
            $("#updateAtHumanize").parent(".ribbon").show();
            $("#updateAtHumanize").html(data['updateAtHumanize']);
        }
    })
        .fail(function() {
            $("#errorMessage").show();
            $("#errorMessage").html("Connection lost!");
        });
}


// When document is ready
(function(){
    console.log("PATH:", "/"+path);

    // Show server errors
    if (error != '') {
        console.log("INTERNAL ERROR:", convertHTMLEntity(error));
    }
    
    // During the first load
    let lastSave = 2;
    reloadNote(lastSave);

    // Load note every 2 seconds
    setInterval(function(){
        lastSave = lastSave+1

        // Reload note
        if (lastSave > 1 ) {
            reloadNote(lastSave);
        }

        // If some fail, force hide the loader
        if (lastSave > 2 ) {
            $("#loader").hide();
        }
    }, 2000);
    $('#noteField').bind('input propertychange', function() {
        // console.log(lastSave);
        lastSave = 0;
    });
    

    // Save content update
    let lastKeyPress;
    $('#noteField').bind('input propertychange', function() {
        $("#loader").transition('show').transition('stop all');
        
        clearTimeout(lastKeyPress);
        lastKeyPress = setTimeout(function() {

            $.post( "/api/"+path , { note: $("#noteField").val() } )
                .done(function(data) {
                    // saved
                    //console.log(data);
                    lastSave = 0;
                    if($("#loader").is(":visible")){
                        $("#loader").transition('zoom');
                    }
                })
                .fail(function() {
                    // fail
                });
        }, 500);
    });
})();
