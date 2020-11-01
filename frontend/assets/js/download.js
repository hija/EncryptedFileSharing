function downloadAndDecrypt(){
    var efsconfigurationame = getUrlParameter("fileId");
    var encryptionKey = window.location.hash.substr(1);

    console.log(efsconfigurationame);
    console.log(encryptionKey);

    encryptedDataDictionary = requestEncryptedFile(efsconfigurationame);
    decryptAndDownload(encryptedDataDictionary, encryptionKey);
}

function getUrlParameter(sParam) {
    var sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
};

function requestEncryptedFile(efsconfigurationame){
     // This data is send to the server now
     return $.ajax({
        url: '/api/download',     
        type: 'post',
        data : JSON.stringify({
            EFSConfigurationName: efsconfigurationame
        }),
        dataType: 'json',
        async: false,                   
        success: function(data)         
        {
            return data;
        } 
    }).responseJSON;
}

function convertWordArrayToUint8Array(wordArray) {
    var arrayOfWords = wordArray.hasOwnProperty("words") ? wordArray.words : [];
    var length = wordArray.hasOwnProperty("sigBytes") ? wordArray.sigBytes : arrayOfWords.length * 4;
    var uInt8Array = new Uint8Array(length), index=0, word, i;
    for (i=0; i<length; i++) {
        word = arrayOfWords[i];
        uInt8Array[index++] = word >> 24;
        uInt8Array[index++] = (word >> 16) & 0xff;
        uInt8Array[index++] = (word >> 8) & 0xff;
        uInt8Array[index++] = word & 0xff;
    }
    return uInt8Array;
}

function decryptAndDownload(input, encryptionKey) {

    var decryptedFileName = CryptoJS.AES.decrypt(input.encryptedFileName, encryptionKey).toString(CryptoJS.enc.Utf8);
    
    var decryptedFile = CryptoJS.AES.decrypt(input.encryptedFile, encryptionKey);
    var typedArray = convertWordArrayToUint8Array(decryptedFile);
    var fileDec = new Blob([typedArray]);

    var a = document.createElement("a");
    var url = window.URL.createObjectURL(fileDec);
    a.href = url;
    a.download = decryptedFileName;
    a.click();
    window.URL.revokeObjectURL(url);
}