// Convert and HTML special char to the UTF-8
function convertHTMLEntity(text){
    const span = document.createElement('span');

    return text
    .replace(/&[#A-Za-z0-9]+;/gi, (entity,position,text)=> {
        span.innerHTML = entity;
        return span.innerText;
    });
}


// When document is ready
(function(){
    console.log("PATH:", "/"+path);

    // Show server errors
    if (error != '') {
        console.log("INTERNAL ERROR:", convertHTMLEntity(error));
    }

    // Load note every 2 seconds
    let lastSave;
    setInterval(function(){
        lastSave = setTimeout(function() {
            $.getJSON( "/api/"+path, function( data ) {
                $("#noteField").val(data['content']);
                $("#noteField").prop("disabled", false);
            });
        }, 100);
    }, 2000);
    $('#noteField').keydown(function() {
        clearTimeout(lastSave);
    });
    

    // Save content update
    let lastKeyPress;
    $('#noteField').keydown(function() {
        clearTimeout(lastKeyPress);
        lastKeyPress = setTimeout(function() {

            $.post( "/api/"+path , { note: $("#noteField").val() } )
                .done(function(data) {
                    // saved
                    //console.log(data);
                    $("#updateAtHumanize").hide();
                })
                .fail(function() {
                    // fail
                });
        }, 100);
    });
})();