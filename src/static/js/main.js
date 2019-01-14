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
function reloadNote () {
    $.getJSON( "/api/"+path, function( data ) {
        $("#loader").hide();
        // console.log(lastSave);
        $("#errorMessage").hide();
        $("#updateAtHumanize").html(data['updateAtHumanize']);
        $("#noteField").val(data['content']);
        $("#noteField").prop("disabled", false);
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
    reloadNote();

    // Load note every 2 seconds
    let lastSave = 0;
    setInterval(function(){
        lastSave = lastSave+1

        if (lastSave > 1 ) {
            reloadNote();
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