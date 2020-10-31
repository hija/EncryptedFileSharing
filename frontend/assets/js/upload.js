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

        // We create a json object which contains all important information, i.e. mimetype + filename
        var unencryptedFileData = {
            'fileName': file.name,
            'fileData': reader.result
        };

        // Now we convert this data --> string which will be encrypted
        var unencryptedFileAsJsonString = JSON.stringify(unencryptedFileData);

        // This data is encrypted now
        var wordArray = CryptoJS.lib.WordArray.create(unencryptedFileAsJsonString);
        var encrypted = CryptoJS.AES.encrypt(wordArray, randomKey).ciphertext;
        
        // Encrypted data is converted to string
        var encryptStr = encrypted.toString(CryptoJS.enc.Base64); // TODO: Maybe chose other encoding, so we save space?

        // This data is send to the server now
        $.ajax({                    
            url: 'http://localhost:8080/api/upload',     
            type: 'post',
            data : JSON.stringify({
              encryptedData : encryptStr
            }),
            dataType: 'json',                   
            success: function(data)         
            {
              console.log(data);
            } 
        });
    };
    reader.readAsArrayBuffer(file);
}

$("#file").change(function(e) {
    var file = e.target.files[0];
    encryptFileAndSend(file);
});