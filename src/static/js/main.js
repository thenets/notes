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

    // Load note for the first time
    
    $.getJSON( "/api/"+path, function( data ) {
        $("#noteField").val(data['content']);
        $("#noteField").prop("disabled", false);
    });

    // Save content update
    let wto;
    $('#noteField').keydown(function() {
        clearTimeout(wto);
        wto = setTimeout(function() {

            $.post( "/api/"+path , { note: $("#noteField").val() } )
                .done(function(data) {
                    // saved
                    $("#updateAtHumanize").hide();
                    console.log(data);
                })
                .fail(function() {
                    // fail
                });
        }, 100);
    });
})();