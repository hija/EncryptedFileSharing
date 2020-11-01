var randomKey = '';

function getRandomKey(){
    // We use AES256 --> We need 256bit.
    // The Word Array returns characters. Because one character is 8 bit we didvide the required bits by 8.
    // Additionally we use .toString() because we will use that random key later in the download-url, so a string is better suited.
    return CryptoJS.lib.WordArray.random(256 / 8).toString();
}

function encryptFileAndSend(file) {
    // Generate a new random key
    var randomKey = getRandomKey();
    // Read the file
    var reader = new FileReader();
    reader.onload = () => {

        // This data is encrypted now
        var wordArray = CryptoJS.lib.WordArray.create(reader.result);
        var encryptedFile = CryptoJS.AES.encrypt(wordArray, randomKey);
        
        // Encrypted data is converted to string
        var encryptedFileStr = encryptedFile.toString();
        var encryptedFileNameStr =  CryptoJS.AES.encrypt(file.name, randomKey).toString();

        // This data is send to the server now
        $.ajax({                    
            url: 'http://localhost:8080/api/upload',     
            type: 'post',
            data : JSON.stringify({
              encryptedFile : encryptedFileStr,
              encryptedFileName: encryptedFileNameStr
            }),
            dataType: 'json',                   
            success: function(data)         
            {
              $('#uploadcontainer').fadeOut(1000, function(){
                $('#shareURL').val(window.location.origin + '/download.html?fileId=' + data.EFSConfigurationName + '#' + randomKey);
                $('#sharingcontainer').fadeIn();
              });
            } 
        });
    };
    reader.readAsArrayBuffer(file);
}

$("#file").change(function(e) {
    var file = e.target.files[0];
    $('#uploadcontainertext').innerHTML = '<h1 class="display-3">We are uploading 🚀</h1><p>We are uploading your file to the server. Please wait!</p>';
    encryptFileAndSend(file);
});